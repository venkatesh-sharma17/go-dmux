Sinks as name implies are the destination for the messages.

Source messages are demultiplexed based on Dmux size and sent to sink.

Dmux size represent number of go routines running each Sink concurrently.

Dmux batch_size could be used to batch messages when they hit sink. Dmux implementation of batching is different. Check Batching Tab for more details.



##HTTPSink
This is a generic HTTPSink, that is built over go http lib.
It provides constructs of pre and post hooks.

This is implemented as infinite retry in case of failure.
Once a message (or batch of messages) reach HTTPSink. HTTPSink will not take next message till this is processed.
