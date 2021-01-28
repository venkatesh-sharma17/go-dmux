
Dmux Connections define Source to Sink connection.
One Dmux instance can run multiple dmux connections.

Configuration options change per connection based on nature of Source and Sink.

Logging is a common configuration option across all connections.

The following Connections are supported:

##kafka_http
This the most simplest connection that connects KafkaSource to HttpSink.

It expects KafkaSource to have KeydMessage - (Key + Value).
ModuloHash will happen on Key if distributor type is Hash (default).

HttpSink will create the following url  : http://endpoint/{topic}/{partition}/{key}/{offset}.
The payload will be byte[] value from Kafka value.
**Note** Dmux does no processing on value.

Any non 2xx response code will be considered as failure by HTTPSink and retried.
Currently retry policy is fixed delay retry. You can configure retry_interval using config.

If batch_size is specified for batching the URL will change to:

http://endpoint/{topic}?batch=(partition1,key1,offset1~partition2,key2,offset2..)
The query param batch was added to help make debugging  easier.

Batching in Dmux is implemented differently, To understand check Batching Tab.

Payload in batch is encoded to convert 2d byte[] to 1d byte[], without any serialization.

###Encoding Format
TODO

###Java Lib
Java library exist to help decode : https://github.com/flipkart-incubator/go-dmux/tree/master/java/godmux-tools

You can add the below POM dependency. (TODO) move this to release version

```pom
<dependency>
  <groupId>com.flipkart.godmux.tools</groupId>
	<artifactId>godmux-tools</artifactId>
	<version>1.0-SNAPSHOT</version>
</dependency>
```


##kafka_foxtrot
This is an extension of kafka_http connection with customization specific to ingest to foxtrot.
https://github.com/Flipkart/foxtrot

Across transact we use foxtrot for monitoring and debugging and have client lib which write to Kafka in a specific message format needed by Foxtrot.

**TODO** Add link to Foxtrot ingestion lib used in transact

This connection works on this structure of Kafka message where key = tableName and value is foxtrot payload.

Check https://github.com/Flipkart/foxtrot/wiki/Event-Ingestion to see foxtrot ingestion API.

URL expected to be configured  http://foxtrot.com:10000/foxtrot/v1/document/__KEY_NAME__

__KEY_NAME__ will be replaced by tableName.

batching is also supported and customised to format needed by foxtrot ingestion api.

query params are added to both Single URL and Batch URL which represent topicName, offset, partition to simplify debugging.
