## Drawbacks:

* ~~Dmux does not support batching.~~ Batching support has now been added. 

* Dmux is stateless and does not support constructs that need persistance such as windowing, aggregations (time or pivot based) and sideline. The expectation is for the app to build such constructs using persistance. (**Note**  HttpSink call fails will infinitely retries and producer will get choked as pending_acks threshold would be hit).

* Very Low latent processing, using local rocksDB per processing to avoid Network overhead can not be achieved with go-dmux. Go-dmux helps in predictive scaling for throughput. Low latency is for App API performance tuning to solve for, but will use a DB isolated over network.




## ~~Sideline (Deprecated now, needs rethinking)~~
~~go-dmux is stateless it requires client to manage persistance for sideline support. If sideline is not configured, go-dmux retries infinitely to execute the httpSink api endpoint till it gets at 2xx response code.~~

~~**Note**: Cost of enabling sideline will be 2x load to client boxes. Clients can avoid enabling this feature and alternative build sideline logic inside their application while processing request and provide 2xx response for performance.~~

~~To enable sideline you just need to add new sink config "sideline_after"~~

```sh
sideline_after = 10 // this implies after 10 retries the message will get sidelined
```


~~Sideline works with 2 constructs, which Client must expose via API.~~

* ~~ShouldTheMessageKeyBeSidelined:~~

 ~~This is a GET call that is made to the same endpoint in the same url format as POST call to execute (**Note** the url is the same url for single message  i.e. http://host:port/path/topicName/partition/key/offset).
 This call is made before every httpSink POST call to validate if the given key needs to be prevented from processing and proactively sidelined as the same key of this message with prior offset would have been already and we want to sideline this as well to ensure ordering.~~

 ~~The API works out of client response codes:~~
 * ~~200 ->  means the key should be proactively sidelined and prevented from processing~~
 * ~~204 ->  means the key can be allowed to process.~~

 ~~**Note** All other status codes result in retry to this GET call infinitely and continuous failure will apply back-pressure on the Source to stall processing.~~

* ~~WriteToSideline:~~

 ~~This is a PUT call that is made to the same endpoint in the same url format as POST call to execute (**Note** the url is the same url for single message i.e. http://host:port/path/topicName/partition/key/offset).
 This call is made post httpSink POST call fails, to sideline the messages. Any Non 2xx response is considered failure~~

 ~~**Note** This PUT call will retry infinitely and continuous failure will apply back-presssure on the Source to stall processing.~~
