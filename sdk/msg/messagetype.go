package msg

// MessageType is a dragonite message type
type MessageType byte

const (
	// DataMsgType message
	DataMsgType MessageType = 0
	// CloseMsgType message
	CloseMsgType MessageType = 1
	// AckMsgType message
	AckMsgType MessageType = 2
	// HeartbeatMsgType message
	HeartbeatMsgType MessageType = 3
)
