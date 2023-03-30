package sideline

type CheckMessageSidelineResponse struct {
	SidelineMessage          bool
	Version                  int32
	MessagePresentInSideline bool
}

type KafkaHttpSidelineMeta struct {
	Endpoint string
	Headers  map[string]string
}

type CheckMessageSideline interface {
	CheckMessageSideline(key []byte) ([]byte, error)
	SidelineMessage(msg []byte) SidelineMessageResponse
	InitialisePlugin(conf []byte) error
}

type SidelineMessage struct {
	GroupId           string
	Partition         int32
	EntityId          string
	Offset            int64
	ConsumerGroupName string
	ClusterName       string
	Message           []byte
	Version           int32
	ConnectionType    string
	SidelineMeta      []byte
}

type SidelineMessageResponse struct {
	Success                     bool
	ConcurrentModificationError error
	UnknownError                error
	ErrorMessage                string
}
