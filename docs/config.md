
## Config Details
DMux reads basic config during start-up from conf.json. It has list of dmuxItems.
On startup DMux reads at run time the config for each dmuxItems to create and start connection.

Basic config include:
* name
* dmuxItems
* logging

Sample config file which which DMux will reading during start-up
```json
{
  "name": "myDmux",
  "dmuxItems": [
    {
      "name": "sample_name",
      "disabled": false,
      "connectionType": "kafka_http",
      "connection": {
        "dmux": {
          "size": 250,
          "distributor_type": "Hash",
          "batch_size": 1
        },
        "source": {
          "name": "source_name",
          "zk_path": "zookeeper1:2181,zookeeper2:2181,zookeeper3:2181/zk-path",
          "kafka_version_major": 2,
          "topic": "source_topic",
          "force_restart": false,
          "read_newest": true
        },
        "pending_acks": 1000000,
        "sink": {
          "endpoint": "http://elb/1.0/api/consume",
          "timeout": "10s",
          "retry_interval": "100ms",
          "headers": [
            {
              "name": "X-Client",
              "value": "go-dumx"
            },
            {
              "name": "Content-Type",
              "value": "application/json"
            }
          ],
          "method": "POST"
        }
      }
    }
  ],
  "logging": {
    "type": "file",
    "config": {
      "path": "default.log",
      "enable_debug": false,
      "rotation": {
        "size_in_mb": 256,
        "retention_count": 5,
        "retention_days": 90,
        "compress": true
      }
    }
  }
}
```

#### Config Details:
| Config Key       | Default | Comment        |
| ------------- |:-------------|:-------------|
| name  | NA | The name given for  this dmux instance|
| dmuxItems  | NA | dmuxItems are dmuxConnections each connection has name and connectionType - name is used to refer to its config and connectionType can be kafka_http or kafka_foxtrot|
| dmux.size  | 10 |demultiplex size. If size = 10; 1 Source will connect to 10 sink. Use this to increase throughput until the client box resource is saturated.   |
| dmux.distributor_type  | Hash |Type of distributor other option is RoundRobin   |
| dmux.batch_size  | 1 | make this value > 1 to specify batching  |
| source.name| NA     | consumer_group_name for Kafka consumer. This will be used in zookeeper offset tracking|
| source.zk_path| NA     | kafka zookeeper path|
| source.topic| NA     | kafka topic you want to consume|
| source.force_restart| false     | set to true to reset consumer to consume from start|
| source.read_newest  |  false    | read from head if this value is set, this config will take in effect only if force_restart is true
| source.kafka_version_major  |  int    | set to 2 if the source is a kafka 2.x.x cluster, 1 if the source is a kafka 1.x.x cluster otherwise ignore it for default (0.8.2)
| sink.endpoint| NA     | http endpoint to hit, If connectionType == kafka_http then  url given here will be appended by /{topic}/{partition}/{key}/{offset}. This will be POST call with byte[] in body, if connectionType == kafka_foxtrot then expected url should be http://foxtrot.com:10000/foxtrot/v1/document/__KEY_NAME__  where __KEY_NAME__ is replaced by kafka-key and body will be JSON. Note: if batch_size is >  1 then batching will result in byte[][] payload for kafka_http connection and []json payload for foxtrot connection|
| sink.timeout| 10s     | http roundtrip timeout |
| sink.retry_interval| 100ms     | time interval to sleep before retry if http call failed. Note: go-dmux has no concept of sideline, It will do infinite retries. Client is expected to build sideline if need at the Sink  Application being hit|
| sink.headers| NA  | static headers to be added in http call. Note:  Content-Type:application/octet-stream will be added for POST calls for kafka_http  and application/json for kafka_foxtrot|
| pending_acks| 10000     | No of unordered acks acceptable till go-dmux starts to apply backpressure to the source. Increase this if QPS does not increase on increasing size and you can see Warning Log in go-dmux that you hit this threshold. Cost of increasing this is memory and larger no of records replay when go-dmux crashes.|
| logging.type| NA | can be either `console` or `file`, decides whether log should be written to console or file |
| logging.config| NA | configuration for `console` or `file` logger |

##### Log config

###### Type: console
| Config Key       | Default | Comment        |
| ------------- |:-------------|:-------------|
| logging.config.enable_debug| false    | boolean flag to enable debug logging|

###### Type: file
| Config Key       | Default | Comment        |
| ------------- |:-------------|:-------------|
| logging.config.path| /var/log/go-dmux/default.log     | log file path|
| logging.config.enable_debug| false    | boolean flag to enable debug logging|
| logging.config.rotation.size_in_mb| 256     | log size|
| logging.config.rotation.retention_count| 5     | number of log files to keep, rest are archived |
| logging.config.rotation.retention_days| 90     | number of days to keep log files, rest are archive. |
| logging.config.rotation.compress| true     | will compress log files which are rotated|


**Note** which every condition becomes true first in retention_count and retention_days will apply.
