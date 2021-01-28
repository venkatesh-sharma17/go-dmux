## Tech Stack
* [go1.12](https://golang.org/)  - golang version
* [Shopify/sarama](https://github.com/Shopify/sarama) - sarma kafka library for kafka servers < 1.0 version
* [wvanbergen/kafka](https://github.com/wvanbergen/kafka) - HA kafka consumer (partition rebalancing etcc.) using zookeeper
* [natefinch/lumberjack.v2](https://gopkg.in/natefinch/lumberjack.v2) - Log rotation
* [mod](https://blog.golang.org/using-go-modules) - mod for dependency management


## Package Info
```sh
├── Gopkg.lock
├── Gopkg.toml
├── README.md
├── build.sh
├── conf.json
├── config.go                          #config integration to read config for dmux conneciton.
├── connection                         #connections
│   ├── kafka_foxtrot_conn.go
│   └── kafka_http_conn.go
├── controller.go                 
├── core                               #core dmux logic
│   ├── distribute.go
│   ├── distribute_test.go
│   ├── dmux.go
│   ├── dmux_test.go
│   ├── util.go
│   └── util_test.go
├── default.log
├── docs
│   ├── Architecture.md
│   ├── Batching.md
│   ├── Benefits.md
│   ├── Onboarding.md
│   ├── Tech.md
│   ├── Test_Setup.md
│   ├── config.md
│   ├── connections.md
│   ├── deployment.md
│   ├── img
│   ├── index.md
│   ├── limitations.md
│   └── monitoring.md
├── http                               #http logic
│   ├── http_sink.go
│   ├── http_sink_test.go
│   ├── sample_order_ids.txt
│   └── sink_test.json
├── images
│   ├── go_dmux.jpg
│   ├── go_dmux_component.png
│   └── kafka_lag_dmux_vs_storm.png
├── java                               #java client tools to consume
│   └── godmux-tools
├── kafka                              #kafa client <1.0 logic
│   ├── kafka_source.go
│   ├── kafka_source_test.go
│   └── offset_tracker.go
├── logging                            #logging code
│   └── logging.go
├── main.go                            #bootstrap main function
├── mkdocs.yml
└── vendor
    ├── github.com
    └── gopkg.in

```