package msg

import (
	"errors"
	"fmt"
)

// ParseMessage parses a raw dragonite message
func ParseMessage(raw []byte) (Message, error) {
	if len(raw) < 2 {
		return nil, errors.New("Packet is too short")
	}

	// The 1st byte in raw is the protocol version, 2nd is the msg type
	msgType := MessageType(raw[1])
	switch msgType {
	case DataMsgType:
		m, err := NewDataMessageFromBytes(raw)
		return &m, err
	case CloseMsgType:
		m, err := NewCloseMessageFromBytes(raw)
		return &m, err
	case AckMsgType:
		m, err := NewAckMessageFromBytes(raw)
		return &m, err
	case HeartbeatMsgType:
		m, err := NewHeartbeatMessageFromBytes(raw)
		return &m, err
	default:
		return nil, fmt.Errorf("Unknown Message Type: %v", byte(msgType))
	}
}
