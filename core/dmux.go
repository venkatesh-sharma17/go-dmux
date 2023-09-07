package core

import (
	"encoding/json"
	"errors"
	"github.com/cenkalti/backoff"
	sideline_module "github.com/flipkart-incubator/go-dmux/sideline"
	"log"
	"math"
	"sync"
	"time"
)

// ControlSignal type defined to build constants used in passing ControlSignals
// to DefaultDmux
type ControlSignal uint8

const (
	//Resize is signal used to trigger resize of a Running Channel
	Resize ControlSignal = 1
	//Stop is signenal used to Stop Dmux
	Stop ControlSignal = 2

	//Sucess response code for ControlSignal Action
	Sucess uint8 = 1
	//Failed response code for ControlSignal Action
	Failed uint8 = 2
)

// DmuxConf holds configuration parameters for Dmux
type DmuxConf struct {
	Size            int             `json:"size"`
	SourceQSize     int             `json:"source_queue_size"`
	SinkQSize       int             `json:"sink_queue_size"`
	DistributorType DistributorType `json:"distributor_type"`
	BatchSize       int             `json:"batch_size"`
	Version         int             `json:"version"`
	Sideline        Sideline        `json:"sideline"`
}

// Sideline holds config parameters for sideline
type Sideline struct {
	Retries               int         `json:"retries"`
	SidelineResponseCodes []int       `json:"sidelineResponseCodes"`
	ConsumerGroupName     string      `json:"consumerGroupName"`
	ClusterName           string      `json:"clusterName"`
	ConnectionType        string      `json:"type"`
	SidelineMeta          interface{} `json:"sidelineMeta"`
}

type DmuxOptionalParams struct {
	EnableDebugLog bool
}

// ControlMsg is the struct passed to Dmux control Channel to enable it
// perform Admin operations such as Resize
type ControlMsg struct {
	signal ControlSignal
	meta   interface{}
}

// ResizeMeta is the struct used to define resize value which is used when
// Dmux is resizing
type ResizeMeta struct {
	newSize int
}

// ResponseMsg is used for running Dmux instnace to response to client
type ResponseMsg struct {
	signal ControlSignal
	status uint8
}

// Sink is interface that implements OutputSink of Dmux operation
type Sink interface {
	// Clone method is expected to return instance of Sink. If Sink is Stateless
	// this can return selfRefrence back. If Sink is Stateful, its good idea to
	// create new instnace of Sink.
	Clone() Sink

	// Consume method gets The interface.
	//TODO currently this method does not return error, need to solve for error
	// handling
	Consume(msg interface{}, retries int, sidelineResponseCodes []int) error

	//BatchConsume method is invoked in batch_size is configured
	BatchConsume(msg []interface{}, version int)
}

// Source is interface that implements input Source to the Dmux
type Source interface {
	//Generate method takes output channel to which it writes data. The
	//implementation can write to to this using multiple goroutines
	//This method is not expected to return, its run in a separate goroutine
	Generate(out chan<- interface{})
	//Method used to trigger GracefulStop of Source
	Stop()
	// GetKey Method to get key for a message
	GetKey(msg interface{}) []byte
	// GetPartition Method to get partition for a message
	GetPartition(msg interface{}) int32
	// GetValue Method to get value for a message
	GetValue(msg interface{}) []byte
	// GetOffset Method to get offsets
	GetOffset(msg interface{}) int64
}

// Distributor interface abstracts the Logic to distribute the load from Source
// to Sink. Client can choose to use HashDistributor or RoundRobinDistributor or
// write their own distribution Logic
type Distributor interface {
	//Distribute method take incoming data interface and number of outbound channels
	//to return the index of channel to be selected for Distribution of this message
	Distribute(data interface{}, size int) int
}

// Dmux struct which enables Size based Dmultiplexing for
// Source to Sink connections.
// TODO restrict size to be powers of 2 for better optimization in modulo
type Dmux struct {
	size                   int
	batchSize              int
	sourceQSize, sinkQSize int
	control                chan ControlMsg
	response               chan ResponseMsg
	err                    chan error
	distribute             Distributor
	version                int
	sideline               Sideline
}

const defaultSourceQSize int = 1
const defaultSinkQSize int = 100
const defaultBatchSize int = 1
const defaultVersion int = 1

// GetDmux is public method used to Get instance of a Dmux struct
func GetDmux(conf DmuxConf, d Distributor) *Dmux {
	control := make(chan ControlMsg)
	response := make(chan ResponseMsg)
	err := make(chan error)
	sourceQSize := defaultSourceQSize
	sinkQSize := defaultSinkQSize
	batchSize := defaultBatchSize
	version := defaultVersion

	if conf.SourceQSize > 0 {
		sourceQSize = conf.SourceQSize
	}

	if conf.SinkQSize > 0 {
		sinkQSize = conf.SinkQSize
	}

	if conf.BatchSize > 0 {
		batchSize = conf.BatchSize
	}

	if conf.Version > 0 {
		version = conf.Version
	}

	output := &Dmux{conf.Size, batchSize, sourceQSize, sinkQSize,
		control, response, err, d, version, conf.Sideline}
	return output
}

//Connect method holds Dmux logic used to Connect Source to Sink
/*func (d *Dmux) Connect(source Source, sink Sink) {
	go d.run(source, sink)
}*/

// Connect method holds Dmux logic used to Connect Source to Sink With Sideline
func (d *Dmux) ConnectWithSideline(source Source, sink Sink, sidelineImpl sideline_module.CheckMessageSideline, optionalParams DmuxOptionalParams) {
	go d.runWithSideline(source, sink, sidelineImpl, optionalParams)
}

// Await method added to enable testing when using bounded source
func (d *Dmux) Await(duration time.Duration) {

	select {
	case e := <-d.err:
		if e != nil {
			panic(e)
		}
	case <-time.After(duration):
		log.Println("Timedout")
	}

}

// Join used to sleep the main routine forever
func (d *Dmux) Join() {
	e := <-d.err
	if e != nil {
		panic(e)
	}
}

// Resize method is used to Resize a running Dmux
func (d *Dmux) Resize(size int) {
	d.control <- getResizeMsg(size)
	<-d.response
}

// Stop is used to GracefulStop running Dmux
func (d *Dmux) Stop() {
	d.control <- getStopMsg()
	<-d.response
}

func getResizeMsg(size int) ControlMsg {
	c := ControlMsg{
		signal: Resize,
		meta:   ResizeMeta{size},
	}
	return c
}

func getStopMsg() ControlMsg {
	c := ControlMsg{
		signal: Stop,
	}
	return c
}

func (d *Dmux) runWithSideline(source Source, sink Sink, sidelineImpl sideline_module.CheckMessageSideline, optionalParams DmuxOptionalParams) {

	ch, wg := setupWithSideline(d.size, d.sinkQSize, d.batchSize, sink, source, d.version, d.sideline, sidelineImpl)
	in := make(chan interface{}, d.sourceQSize)
	//start source
	//TODO handle panic recovery if in channel is closed for shutdown
	go source.Generate(in)

	for {
		select {
		case data := <-in:
			i := d.distribute.Distribute(data, len(ch))
			// log.Printf("writing to channel %d len %d", i, len(ch[i]))
			if optionalParams.EnableDebugLog {
				log.Printf("writing to channel %d len %d", i, len(ch[i]))
			}
			ch[i] <- data
		case ctrl := <-d.control:
			if ctrl.signal == Resize {
				log.Println("processing resize")
				shutdown(ch, wg)
				resizeMeta := ctrl.meta.(ResizeMeta)
				ch, wg = setupWithSideline(resizeMeta.newSize, d.sinkQSize, d.batchSize, sink, source, d.version, d.sideline, sidelineImpl)
				d.response <- ResponseMsg{ctrl.signal, Sucess}
			} else if ctrl.signal == Stop {
				log.Println("processing stop")
				source.Stop()
				shutdown(ch, wg)
				close(in)
				d.response <- ResponseMsg{ctrl.signal, Sucess}
				d.err <- nil
				return
			}
		}
	}
}

func shutdown(ch []chan interface{}, wg *sync.WaitGroup) {
	for _, c := range ch {
		close(c)
	}
	wg.Wait()
}

/*
	func setup(size, qsize, batchSize int, sink Sink, version int) ([]chan interface{}, *sync.WaitGroup) {
		if version == 1 && batchSize == 1 {
			return setupWithSideline(size, qsize, batchSize, sink, version, nil, nil)
		} else {
			return batchSetup(size, qsize, batchSize, sink, version)
		}
	}
*/
func setupWithSideline(size, qsize, batchSize int, sink Sink, source Source, version int, sideline Sideline, sidelineImpl sideline_module.CheckMessageSideline) ([]chan interface{}, *sync.WaitGroup) {
	if version == 1 && batchSize == 1 {
		if sidelineImpl != nil {
			log.Printf("Calling simpleSetupWithSideline")
			return simpleSetupWithSideline(size, qsize, sink, source, sideline, sidelineImpl)
		} else {
			log.Printf("Calling simpleSetup")
			return simpleSetup(size, qsize, sink)
		}
	} else {
		if sidelineImpl == nil {
			return batchSetup(size, qsize, batchSize, sink, version)
		}
		log.Fatal("Not Supported sidelining for batching")
		return nil, nil
	}
}

// Batching has been implemented a bit differently to provide more value to end Client.
// Typical batching would involve queing up till batch size and flush.
// This is inefficient at Client consumption when using modulo hashing (HashDistributor).
// Client of Batch Sink has to run single threaded code to ensure ordering is not lost at consumption.
// The implementation below avoids this shortcoming.
// In this implementation we create total = size * batchSize number of channels.
// And we create numer of BatchConsumer routine = size.
// Each BatchConsumer routine will consume (index, index+batchSize); index = index of consumer
// which is 0 to size-1.
// BatchConsumer will update its batch array index from one entry each of respective channel index. (This provides
// ability for consumer to consume in parallel) and then flush the batch.
// Close of any channel in a BatchConsumer will stop the BatchConsumer.
func batchSetup(sz, qsz, batchsz int, sink Sink, version int) ([]chan interface{}, *sync.WaitGroup) {
	size := sz * batchsz // create double nuber of channels

	wg := new(sync.WaitGroup)
	wg.Add(size)
	ch := make([]chan interface{}, size)

	//create channels
	for i := 0; i < size; i++ {
		ch[i] = make(chan interface{}, qsz)
	}

	//assign batch of channels per batch consumer
	for i := 0; i < size; i += batchsz {
		//async runner does batching and call to sink.BatchConsume
		go func(index int) {
			sk := sink.Clone()
			batch := make([]interface{}, batchsz) //reuse one time created array

			for { // run till any channel closes check failed case which breaks from this loop
				// log.Println("Running consumer ", index)
				//batch
				failed := false
				j := index
				for z := 0; z < batchsz; z++ {
					// log.Printf("Iterating count %d consumer %d j %d channelLen %d \n", z, index, j, len(ch[j]))
					msg, more := <-ch[j]  // more == false if channel is closed
					if !failed && !more { // failed = true if any one channel is closed
						failed = true
						break
					}
					batch[z] = msg
					j++
				}
				// log.Println("Failed ", failed)
				if failed { //ack all waitGroups and return to break from the infinite for loop
					for z := 0; z < batchsz; z++ {
						wg.Done()
					}
					return
				}
				// log.Println("flusing ", batch)
				//flush batched message
				sk.BatchConsume(batch, version)
			}

		}(i)
	}
	return ch, wg
}

func simpleSetup(size, qsize int, sink Sink) ([]chan interface{}, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	wg.Add(size)
	var responseCodes []int
	ch := make([]chan interface{}, size)
	for i := 0; i < size; i++ {
		ch[i] = make(chan interface{}, qsize)
		go func(index int) {
			sk := sink.Clone()
			for msg := range ch[index] {
				sk.Consume(msg, math.MaxInt32, responseCodes)
			}
			wg.Done()
		}(i)
	}
	return ch, wg
}

func sinkConsume(sink Sink, sinkChannel []chan ChannelObject, index int, sideline Sideline, sidelineChannel []chan ChannelObject) {
	sk := sink.Clone()
	expBackOff := backoff.NewExponentialBackOff()
	//expBackOff.MaxElapsedTime = math.MaxInt32 * time.Minute
	for channelObject := range sinkChannel[index] {
		log.Printf("Inside Sink channel ")
		retryError := backoff.Retry(func() error {
			consumeError := sk.Consume(channelObject.Msg, sideline.Retries, sideline.SidelineResponseCodes)
			if consumeError == nil {
				return nil
			}
			if consumeError != nil && consumeError.Error() == SidelineMessage {
				sendToSidelineChannel := ChannelObject{
					Msg:      channelObject.Msg,
					Sideline: channelObject.Sideline,
					Version:  0,
				}
				sidelineChannel[index] <- sendToSidelineChannel
				return nil
			}
			return errors.New("failed in sink consume " + consumeError.Error())
		}, expBackOff)
		if retryError != nil {
			log.Fatal("Ideally this should not happen in sinkConsume" + retryError.Error())
		}
	}
}

func mainChannelConsumption(ch []chan interface{}, index int, source Source, sideline Sideline, sidelineImpl sideline_module.CheckMessageSideline,
	sidelineChannel []chan ChannelObject, sinkChannel []chan ChannelObject, wg *sync.WaitGroup) {
	for msg := range ch[index] {
		key := source.GetKey(msg)
		partition := source.GetPartition(msg)
		value := source.GetValue(msg)
		offset := source.GetOffset(msg)
		var check sideline_module.CheckMessageSidelineResponse
		expBackOff := backoff.NewExponentialBackOff()
		//expBackOff.MaxElapsedTime = math.MaxInt32 * time.Minute
		retryError := backoff.Retry(func() error {
			log.Printf("Checking if the message is already sidelined %d, %d from channel", partition, offset)
			checkSidelineMessage := sideline_module.SidelineMessage{
				GroupId:           string(key),
				Partition:         partition,
				EntityId:          string(key) + sideline.ConsumerGroupName + sideline.ClusterName,
				Offset:            offset,
				ConsumerGroupName: sideline.ConsumerGroupName,
				ClusterName:       sideline.ClusterName,
				Message:           value,
				ConnectionType:    sideline.ConnectionType,
			}
			checkSidelineMessageBytes, checkSidelineMessageErr := json.Marshal(checkSidelineMessage)
			if checkSidelineMessageErr != nil {
				log.Printf("error in serde of checkSidelineMessage " + checkSidelineMessageErr.Error())
				return errors.New("error in serde of checkSidelineMessage " + checkSidelineMessageErr.Error())
			}
			checkBytes, checkErr := sidelineImpl.CheckMessageSideline(checkSidelineMessageBytes)
			err := json.Unmarshal(checkBytes, &check)
			if err != nil {
				log.Printf("error in serde of CheckMessageSidelineResponse " + err.Error())
				return errors.New("error in serde of CheckMessageSidelineResponse " + err.Error())
			}
			if checkErr != nil {
				log.Printf("Error in checking if message is sidelined " + checkErr.Error())
				return errors.New("Error in checking if message is sidelined " + checkErr.Error())
			}
			log.Printf("Message if already sidelined %t %d %d", check.MessagePresentInSideline, partition, offset)
			if check.MessagePresentInSideline {
				return nil
			}
			log.Printf("SidelineMessage %t %d %d", check.SidelineMessage, partition, offset)
			if check.SidelineMessage {
				sendToSidelineChannel := ChannelObject{
					Msg:      msg,
					Sideline: sideline,
					Version:  check.Version,
				}
				sidelineChannel[index] <- sendToSidelineChannel
				return nil
			} else {
				sendToSinkChannel := ChannelObject{
					Msg:      msg,
					Sideline: sideline,
				}
				sinkChannel[index] <- sendToSinkChannel
				return nil
			}
		}, expBackOff)
		if retryError != nil {
			log.Fatal("Ideally this should not happen in mainChannelConsumption" + retryError.Error())
		}
	}
	wg.Done()
}

func pushToSideline(sidelineChannel []chan ChannelObject, index int, source Source, sideline Sideline, sidelineMetaByteArray []byte, sidelineImpl sideline_module.CheckMessageSideline) {
	for channelObject := range sidelineChannel[index] {
		expBackOff := backoff.NewExponentialBackOff()
		//expBackOff.MaxElapsedTime = math.MaxInt32 * time.Minute
		retryError := backoff.Retry(
			func() error {
				val := source.GetValue(channelObject.Msg)
				key := source.GetKey(channelObject.Msg)
				partition := source.GetPartition(channelObject.Msg)
				offset := source.GetOffset(channelObject.Msg)
				log.Printf("Inside sideline channel for partition %d offset %d", partition, offset)
				kafkaSidelineMessage := sideline_module.SidelineMessage{
					GroupId:           string(key),
					Partition:         partition,
					EntityId:          string(key) + sideline.ConsumerGroupName + sideline.ClusterName,
					Offset:            offset,
					ConsumerGroupName: sideline.ConsumerGroupName,
					ClusterName:       sideline.ClusterName,
					Message:           val,
					Version:           channelObject.Version,
					ConnectionType:    sideline.ConnectionType,
					SidelineMeta:      sidelineMetaByteArray,
				}

				sidelineByteArray, err := json.Marshal(kafkaSidelineMessage)
				if err != nil {
					log.Printf("error in serde of kafkaSidelineMessage " + err.Error())
					return errors.New("error in serde of kafkaSidelineMessage " + err.Error())
				}
				sidelineMessageResponse := sidelineImpl.SidelineMessage(sidelineByteArray)
				if !sidelineMessageResponse.Success {
					var check sideline_module.CheckMessageSidelineResponse
					log.Printf(sidelineMessageResponse.ErrorMessage)
					if sidelineMessageResponse.ConcurrentModificationError != nil {
						checkSidelineMessage := sideline_module.SidelineMessage{
							GroupId:           string(key),
							Partition:         partition,
							EntityId:          string(key) + sideline.ConsumerGroupName + sideline.ClusterName,
							Offset:            offset,
							ConsumerGroupName: sideline.ConsumerGroupName,
							ClusterName:       sideline.ClusterName,
							Message:           val,
							ConnectionType:    sideline.ConnectionType,
						}
						checkSidelineMessageBytes, checkSidelineMessageErr := json.Marshal(checkSidelineMessage)
						if checkSidelineMessageErr != nil {
							log.Printf("error in serde of checkSidelineMessage")
							return errors.New("error in serde of checkSidelineMessage")
						}
						checkBytes, checkErr := sidelineImpl.(sideline_module.CheckMessageSideline).
							CheckMessageSideline(checkSidelineMessageBytes)
						err := json.Unmarshal(checkBytes, &check)
						if err != nil {
							log.Printf("error in serde of CheckMessageSidelineResponse")
							return errors.New("error in serde of CheckMessageSidelineResponse")

						}
						if checkErr != nil {
							log.Printf("Error in checking if message is sidelined " + checkErr.Error())
							return errors.New("Error in checking if message is sidelined " + checkErr.Error())
						}
						channelObject.Version = check.Version
						return errors.New("Retrying ")
					}
					return errors.New("Retrying ")
				}
				return nil
			}, expBackOff)
		if retryError != nil {
			log.Fatal("Ideally this should not happen in pushToSideline")
		}
	}
}

func simpleSetupWithSideline(size, qsize int, sink Sink, source Source, sideline Sideline, sidelineImpl sideline_module.CheckMessageSideline) ([]chan interface{}, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	wg.Add(size)
	ch := make([]chan interface{}, size)
	sidelineChannel := make([]chan ChannelObject, size)
	sinkChannel := make([]chan ChannelObject, size)
	log.Printf("Inside simpleSetupWithSideline")
	for i := 0; i < size; i++ {
		sinkChannel[i] = make(chan ChannelObject, qsize)
		go sinkConsume(sink, sinkChannel, i, sideline, sidelineChannel)
	}

	for i := 0; i < size; i++ {
		sidelineChannel[i] = make(chan ChannelObject, qsize)
		sidelineMetaByteArray, sidelineMetaByteArrayErr := json.Marshal(sideline.SidelineMeta)
		if sidelineMetaByteArrayErr != nil {
			log.Fatal("error in serde of SidelineMeta")
		}
		go pushToSideline(sidelineChannel, i, source, sideline, sidelineMetaByteArray, sidelineImpl)
	}

	for i := 0; i < size; i++ {
		ch[i] = make(chan interface{}, qsize)
		go mainChannelConsumption(ch, i, source, sideline, sidelineImpl, sidelineChannel, sinkChannel, wg)
	}
	return ch, wg
}
