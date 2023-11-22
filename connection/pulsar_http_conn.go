package connection

import (
	"encoding/json"
	"github.com/flipkart-incubator/go-dmux/core"
	sink "github.com/flipkart-incubator/go-dmux/http"
	source "github.com/flipkart-incubator/go-dmux/pulsar"
	"log"
)

// PulsarConnConfig holds config to connect pulsar source to http sink
type PulsarConnConfig struct {
	Dmux        core.DmuxConf     `json:"dmux"`
	Source      source.PulsarConf `json:"source"`
	Sink        sink.HTTPSinkConf `json:"sink"`
	PendingAcks int               `json:"pending_acks"`
}

// PulsarConn abstracts connection
type PulsarConn struct {
	EnableDebugLog bool
	Conf           interface{}
}

// getConfiguration parses configs and returns connection config
func (c *PulsarConn) getConfiguration() *PulsarConnConfig {
	data, _ := json.Marshal(c.Conf)
	var config *PulsarConnConfig
	json.Unmarshal(data, &config)
	return config
}

// Run starts connection from source to sink
func (c *PulsarConn) Run() {
	conf := c.getConfiguration()
	log.Println("starting go-dmux with conf", conf)

	src := source.GetPulsarSource(conf.Source)
	tracker := source.GetCursorTracker(conf.PendingAcks, src)
	hook := source.GetPulsarHook(tracker, c.EnableDebugLog)

	snk := sink.GetHTTPSink(conf.Dmux.Size, conf.Sink)
	snk.RegisterHook(hook)
	src.RegisterHook(hook)

	h := source.GetMessageHasher()
	d := core.GetDistribution(conf.Dmux.DistributorType, h)

	dmux := core.GetDmux(conf.Dmux, d)
	var optionalParams core.DmuxOptionalParams = core.DmuxOptionalParams{c.EnableDebugLog}
	dmux.ConnectWithSideline(src, snk, nil, optionalParams)
	dmux.Join()
}
