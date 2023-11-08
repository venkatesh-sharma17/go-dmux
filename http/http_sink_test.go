package http

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type GetOrderMsg struct {
	orderID string
}

type OrderMsgHasher struct{}

func (o *OrderMsgHasher) ComputeHash(data interface{}) int {
	val := data.(GetOrderMsg)
	h := fnv.New32a()
	h.Write([]byte(val.orderID))
	return int(h.Sum32())
}

func (g GetOrderMsg) GetMethod(conf HTTPSinkConf) string {
	return "GET"
}

func (g GetOrderMsg) GetPayload() []byte {
	return nil
}

func (g GetOrderMsg) IsSidelined() bool {
	return false
}

func (g GetOrderMsg) Sideline() {

}

func (g GetOrderMsg) GetHeaders(conf HTTPSinkConf) map[string]string {
	header := map[string]string{
		"X_CLIENT_ID": "TEST",
	}
	return header
}

func (g GetOrderMsg) GetURLPath() string {
	return "/path/" + g.orderID
}

type PrintHook struct{}

func (p *PrintHook) PreHTTPCall(msg interface{}) {
	// do nothing
}
func (p *PrintHook) PostHTTPCall(msg interface{}, success bool) {
	log.println(msg, success)
}

type FileSource struct {
	inputFilePath string
}

func GetFileSource(path string) *FileSource {
	return &FileSource{path}
}

func (f *FileSource) Generate(out chan<- interface{}) {

	ordersBytes, err := ioutil.ReadFile(f.inputFilePath)
	if err != nil {
		log.Fatalf("Failed to read file %s, %v", f.inputFilePath, err)
		return
	}

	ordersString := strings.TrimSpace(string(ordersBytes))
	orders := strings.Split(ordersString, "\n")
	for _, orderID := range orders {
		msg := GetOrderMsg{orderID}
		out <- msg
	}

}

func (f *FileSource) Stop() {

}

func parseConf(path string) HTTPSinkConf {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	var conf HTTPSinkConf
	json.Unmarshal(raw, &conf)

	return conf
}

// func TestHTTPSink(t *testing.T) {
// 	log.println("running test TestHTTPSink")
// 	hook := new(PrintHook)
// 	conf := parseConf("sink_test.json")
// 	sink := GetHTTPSink(10, conf)
// 	sink.RegisterHook(hook)
// 	msg := GetOrderMsg{"123"}
// 	sink.Consume(msg)
// }

// func TestHTTPSinkWithDMux(t *testing.T) {
// 	log.println("running test TestHTTPSinkWithDMux")
// 	hasher := new(OrderMsgHasher)
// 	d := core.GetHashDistribution(hasher)
// 	dconf := core.DmuxConf{
// 		Size:        4,
// 		SourceQSize: 1,
// 		SinkQSize:   100,
// 	}
// 	dmux := core.GetDmux(dconf, d)
// 	hook := new(PrintHook)
// 	source := GetFileSource("sample_order_ids.txt")
// 	conf := parseConf("sink_test.json")
// 	sink := GetHTTPSink(dconf.Size, conf)
// 	sink.RegisterHook(hook)
//
// 	dmux.Connect(source, sink)
// 	dmux.Await(1 * time.Second)
// 	dmux.Resize(10)
// 	dmux.Await(1 * time.Second)
// }
