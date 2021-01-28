
## Test Setup:
* Number of Storm Workers == Number of App Boxes behind VIP == Partitions of this Kafka topic == 8.
* AppBoxes running Dropwizard with 256 worker threads. The boxes were behind VIP which go-dmux was pointing to.
* Single go-dmux instance with size = 2000, pending_acks = 50000.
* Strom worker tried with max_spout_pending = 100000, and executors count = 256, 512, 1024.
* Logic in App == Logic in Spout + Bolt.
* Both storm and dmux where run with force_restart = true.
* We tried restarting Storm with different configuration to enable better performance after go-dmux drained. Hence multiple peaks are present in the graph.


![](./../img/kafka_lag_dmux_vs_storm.png)

## We believe better Performance of Dmux is because:
* Dmux avoids cross JVM chatter for single request processing. In Storm there are multiple queue hops in Spout -> Bolt.
* Dmux does not have any single choke point for its processing. In Storm Zk can be bottleneck as you increase number of executors and workers.
* Dmux has much simpler offset tracking which runs within single Dmux routine for the source partitions it owns. There is not cross process chatter for offset tracking. In Strom cross chatter leads to un-predictable performance when increasing max_spout_pending.

