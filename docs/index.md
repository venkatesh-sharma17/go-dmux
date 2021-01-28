
# About
A Simple Stateless service which can DMux Source to Sink.
This was Built for High Throughput Processing over KafkaStream Source and HttpEndpoint Sink.

## Application Overview
Core of go-dmux is to take any source and connect to any sink. Dmux takes distributor and size as argument which enables it to distribute data from source to sink in a FanOut manner.

KafkaSource and HttpSink have been added by default, as that was the major use case Dmux was built to solve for. Prior to Dmux we were using Storm and Spark for Stream processing applications. Dmux is a simpler and predictive alternative to Storm, Spark or other Stream processing systems. Dmux goal is have developer focus on writing app's and have go-dmux to redirect messages from Kafka to app's. This enables reuse of App monitoring infrastructure which is already in place and also in NFR tuning.

The key difference between Dmux and Storm is much lesser number of hops and queues involved in Single Message processing. Dmux has single hop from Source to Sink. With reference to [queueing theory](https://en.wikipedia.org/wiki/Queueing_theory#Queueing_networks), ability to estimate when a system will be queued up and degrade is simpler with Dmux (M/M/1). This gets fairly complicated with Storm which has two way queues per Spout, Bolt and OffsetAckers. Dmux helps in making it more simple to predict the number boxes needed to scale.

**Note** Go-Dmux does not support the concepts around higher order processing primitives such as ParDo or Aggregate by time window. These are not primitives in Storm as well. As Trident is Built over Storm such constructs can be built over App Framework using Datastore.
