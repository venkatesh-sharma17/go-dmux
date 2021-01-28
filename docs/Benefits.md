
### Benefits

#### Tech Stack Consolidation (Better Reuse and Developer Productivity):
* Developers focus only on Building API's over Apps irrespective of Stream Processing or User request handling. No need to worry about new components such as Nimbus, Supervisor in case of Storm.
* Monitoring and Scaling are solved problems for HA Apps. Developers can reuse patterns of dispatch compose and  state machine orchestration used in Apps.
* Single Code base for App and Storm makes it easier to maintain. Just add another API to build stream processing logic.
* App Developers do not need to worry about processing message in order. Dmux provides ordering per Key guarantees by ensuring the next call for the same key is only made when first call returns 2xx.


#### Simpler Predictive Scalability:
* Simple way to scale with Dmux is by increasing fanout concurrency or increasing number of App boxes acting as HTTPSink or both depending on which is the bottleneck.

For eg.

**Note:** This is just a simple approach which works for most cases.
```
If avg latency for processing a message is 500ms. (Very Slow processing use case)
And Single App can process 200 concurrent request at 70% CPU saturation.
Then we can get Throughput of 400 req/sec from a single box.
Now we can meet given QPS needs by just adding the right number of boxes and increasing Dmux concurrency.

For meeting 2000 QPS.
 - We need at least 5 App boxes behind the HTTPSink VIP.  = 5 * 400
 - We need increase fanout of size of Dmux = 1000 concurrent request = 5 * 200
```
* We can achieve Vertical Scalability in single Dmux instance by touching two properties in Dmux config, **size** and **pending_acks**. Increasing them till we see good utilization of VM

 **Note:** you are more likely to saturate SinkApp boxes before you saturate at a single Dmux instance

* We can achieve Horizontal Scalability by increasing number of boxes for Dmux or sinkApp boxes, based on which is saturated.

 **Note:** Dmux horizontal scalability limit is limited by number of Kafka Partitions in case of KafkaSource. It's unlikely you will hit this limit given Dmux's vertical scale.

 **Update** Batching is now added, which helps in better throughput.

#### Easier to achieve Better Performance
Dmux performance was much better than Storm-0.9.5. Our basic testing had go-dmux draining Kafka in 1/10 time of Storm-0.9.5. Storm-0.9.5 took close to 5hrs to drain the queue when go-dmux could do it in 30min.

**Note**: This was with our default tuning of Storm-0.9.5. We suspect there was backpreassure being built up and our basic tuning with count of workers, executors, max_spout_pending did not help. We did not dig deeper into Storm to find the bottleneck and solve for it in this testing.

Check Graph in TestSetup which displays Kafka Lag monitoring Graph that shows Drain rate of same Topology in Storm vs Dmux + HttpAPP end point when the KafkaConsumer when force restarted.
