package main

import (
	co "github.com/flipkart-incubator/go-dmux/config"
	"github.com/flipkart-incubator/go-dmux/metrics"
	"log"
	"net/http"
	"os"

	"github.com/flipkart-incubator/go-dmux/logging"
)

import _ "net/http/pprof"

//

// **************** Bootstrap ***********

func main() {
	args := os.Args[1:]
	sz := len(args)

	var path string

	if sz == 1 {
		path = args[0]
	}
	log.Println("preparing go routines for pprof")
	go func() {
		log.Println("starting web server for pprof")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	dconf := co.DMuxConfigSetting{
		FilePath: path,
	}
	conf := dconf.GetDmuxConf()

	dmuxLogging := new(logging.DMuxLogging)
	dmuxLogging.Start(conf.Logging)

	log.Printf("config: %v \n", conf)

	//start showing metrics at the endpoint
	metrics.Start(conf.MetricPort)

	for _, item := range conf.DMuxItems {
		go func(connType co.ConnectionType, connConf interface{}, logDebug bool) {
			connType.Start(connConf, logDebug, nil)
		}(item.ConnType, item.Connection, dmuxLogging.EnableDebug)
	}

	//main thread halts. TODO make changes to listen to kill and reboot
	select {}
}
