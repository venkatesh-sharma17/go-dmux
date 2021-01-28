
# go-dmux

# About
A Simple Stateless service which can DMux Source to Sink.
This was Built for High Throughput Processing over KafkaStream Source and HttpEndpoint Sink.

## Features
* DMux any Source to any Sink
* High throughput processing
* Ordering guarantees
* Batching capability
* Two distribution type supported (RoundRobin and Hash)
* Fork based on connection type, currently supporting kafka to Http and kafka to foxtrot
* single  Dmux for multiconsumer support
* Rich unit tests
* Rich documentatation

**Note** Go-Dmux does not support the concepts around higher order processing primitives such as ParDo or Aggregate by time window. These are not primitives in Storm as well. As Trident is Built over Storm such constructs can be built over App Framework using Datastore.

## License
Copyright 2021 Flipkart Internet, pvt ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## Documentation
Refer [here](https://github.com/pages/flipkart-incubator/go-dmux/) for documentation.