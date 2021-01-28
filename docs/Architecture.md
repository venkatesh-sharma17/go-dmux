### Dmux Components

![svg](./../img/go_dmux_blocks.png)

* DmuxConnection connects one Source to Sink.
* DmuxInstance can run multiple DmuxConnections.
* Multiple DmuxInstances running same DmuxConnection will load balance.
* LoadBalance is function of Source. In case of KafkaSource it load balances using Zookepeer.
* DmuxInstance is HA stateless component which you can install on any compute VM and scale out to the number of Kafka Partitions. Good to use Instance group to resize.


### DmuxConnection Block Digram


![](./../img/go_dmux.jpg)

* DmuxConnection connects Source to Sink. Example connectionTypes = kafka_http, kafka_foxtrot.
* DmuxConnection enables you to configure dmuxSize concurrent consumer Sinks along with batching to achieve high Throughput.
* Dmux is written Go and leverages go coroutines and channels to achieve high vertical scale. Most cases would need just one Dmux instance to run many DmuxConnections.
* KafkaSource interface in Dmux is a High Available KafkaConsumer which reuses the Zookeeper used by KafkaBrokers for Partition Balancing and Offset management.
* During Partition Rebalancing between Dmux instances, Client can expect replay but there will be no data loss.
* Dmux ensures ordering per Key when distributor = Hash. Note this is default distributor type.
