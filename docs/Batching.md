#Batching
Batching capability enables you to specify **batch_size**. Batching logic is implemented differently in go-dmux to help achieve better performance. When using Modulo Hash Distributor, Batching is normally implemented by request collapsing per Consumer.
Dmux implements batching by increasing number of channels = batchSize * size and having one consumer per batch of channels. Consumption inside a consumer is to pop the head from each channel serially and build batch request for sink. This ensures batch entities are not of same group (Hash group) and can be processed in parallel by sink.

Given below is an illustration of the above idea.
```
eg.. if Dmux size == 2 and batch_size == 2;
 source = [O1t1,O2t1,O1t2,O1t3,O3t1,O3t2,O2t2,O4t1,O5t1,O3t3]
 # O1 - Order1 and t1 = time 1
 # hence O1t1 = Order 1 at time 1
  -> Lets assume this is distributed using HashDistributor as:
  sink1 - [O1t1,O1t2,O1t3,O4t1,O5t1]
  sink2 - [O2t1,O3t1,O3t2,O2t2,O3t3]

  Now Traditional batching would batch this before calling Sink as
  sink1 - [[O1t1,O1,t2], [O1t3,O4,t1], [O5t1]]
  similar for sink2.
  The problem here is that when the App endpoint processing this HTTPSink call to it needs to process this, it has to process sequentially to ensure ordering is maintained.

  Imaging Batch 1 from Sink 1 = [O1t1, O1t2] if processed parallel has potential list of out of order in processing.

  Dmux implementation of Batching for same source :
  Dmux creates channels = batchsize * dmuxSize => 2* 2 = 4
 -> Now the  HashDistributor will distribute this to the 4 channels as

  channel1 = [O1t1,O1t2,O1t3]
  channel2 = [O2t1,O2t2]
  channel3 = [O3t1,O3t2,O3t3]
  channel4 = [O4t1,O5t1]

  Sink to channel mapping would be
  sink1 = [channel1, channel2]
  sink2 = [channel3, channel4]

  each sink ensures it batches by reading top entry from each of the channel.
  Hence batch will contain entries which are not of the same groupId. Enabling client o process them concurrently.
```
