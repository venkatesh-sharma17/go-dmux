package sideline_impls

import (
	co "github.com/flipkart-incubator/go-dmux/config"
	"github.com/flipkart-incubator/go-dmux/logging"
	"github.com/flipkart-incubator/go-dmux/metrics"
	"log"
)

//

type DmuxCustom struct {
}

func (d *DmuxCustom) DmuxStart(path string, sidelineImp interface{}) {
	//log.Println(checkMessageSideline.SidelineMessage())

	dconf := co.DMuxConfigSetting{
		FilePath: path,
	}
	conf := dconf.GetDmuxConf()

	dmuxLogging := new(logging.DMuxLogging)
	//_ = new(logging.DMuxLogging)

	log.Printf("config: %v", conf)

	//start showing metrics at the endpoint
	metrics.Start(conf.MetricPort)

	for _, item := range conf.DMuxItems {
		log.Println(item.ConnType)
		go func(connType co.ConnectionType, connConf interface{}, logDebug bool) {
			connType.Start(connConf, logDebug, sidelineImp)
		}(item.ConnType, item.Connection, dmuxLogging.EnableDebug)
	}

	//main thread halts. TODO make changes to listen to kill and reboot
	select {}
}
