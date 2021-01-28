Sources implement logic go generate infinite stream of messages.
Sources also have capability of load balancing when go-dmux connection is running with multiple go-dmux instance to enable distribution of incoming stream between the multiple go-dmux instances running the same connection.


##KafaSource

Kafka source is wrapper implementation over wvanbergen/kafka/consumergroup which uses sarama.kafka. sarama.kafka was chosen as it supported kafka 0.8 version. wvanbergen/kafka/consumergroup provide HA worker constructs to handle offset using zookeeper and load balance partition of kafka topic across go-dmux instances having the same connection.


Kafka Consumer created by KafkaSource is similar to KafkaConsoleConsumer, it creates entries under consumer node of kafkaPath in zookeeper that holds the consumer offset.
Lag in topic consumption can be monitored by checking offset diff between producer and consumer.

For Transact team, a logic to collect this metric and push to JMX exist and can be configured.

TODO - add notes on how to configure

Other team can setup the daemon process to collect and push metric to JMX.

Project with code to monitor Kafka offsets via Zookeeper:

TODO - (move above project to transact-commons add documentation to help setup monitoring)

###Note:
Always setup monitoring and track lag, You should also setup alerting on this lag. This is the first level indicator when things don't work
