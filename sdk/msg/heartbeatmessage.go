package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dragonite-network/go-dragonite/sdk/misc"
)

// HeartbeatMessage is a heartbeat message
type HeartbeatMessage struct {
	seq int32
}

// NewHeartbeatMessage constructs a heartbeat message from a sequence
func NewHeartbeatMessage(seq int32) HeartbeatMessage {
	return HeartbeatMessage{
		seq,
	}
}

// NewHeartbeatMessageFromBytes constructs a heartbeat message from raw msg bytes
func NewHeartbeatMessageFromBytes(raw []byte) (HeartbeatMessage, error) {
	var remoteVersion byte
	var remoteType byte
	var seq int32

	buff := bytes.NewBuffer(raw)

	if err := binary.Read(buff, binary.BigEndian, &remoteVersion); err != nil {
		return HeartbeatMessage{}, err
	}
	if remoteVersion != misc.ProtocolVersion {
		return HeartbeatMessage{}, fmt.Errorf("Incorrect version (%v, should be %v)", remoteVersion, misc.ProtocolVersion)
	}
	if err := binary.Read(buff, binary.BigEndian, &remoteType); err != nil {
		return HeartbeatMessage{}, err
	}
	if remoteType != byte(HeartbeatMsgType) {
		return HeartbeatMessage{}, fmt.Errorf("Incorrect type (%v, should be %v)", remoteType, byte(HeartbeatMsgType))
	}
	if err := binary.Read(buff, binary.BigEndian, &seq); err != nil {
		return HeartbeatMessage{}, err
	}

	return HeartbeatMessage{
		seq,
	}, nil
}

// Version returns the protocol version
func (m *HeartbeatMessage) Version() byte {
	return misc.ProtocolVersion
}

// Type returns message type
func (m *HeartbeatMessage) Type() MessageType {
	return HeartbeatMsgType
}

// FixedLength returns the length of the fixed part of the msg
func (m *HeartbeatMessage) FixedLength() uint16 {
	return 8
}

// Bytes returns the message in raw byte array
func (m *HeartbeatMessage) Bytes() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	if err := binary.Write(buff, binary.BigEndian, byte(m.Version())); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, byte(m.Type())); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, int32(m.seq)); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
