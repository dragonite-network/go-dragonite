package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dragonite-network/go-dragonite/sdk/misc"
)

// CloseMessage is a message signaling close
type CloseMessage struct {
	seq    int32
	status int16
}

// NewCloseMessage constructs a close message from sequence number and status
func NewCloseMessage(seq int32, status int16) CloseMessage {
	return CloseMessage{
		seq,
		status,
	}
}

// NewCloseMessageFromBytes constructs a close message from a raw msg
func NewCloseMessageFromBytes(raw []byte) (CloseMessage, error) {
	var remoteVersion byte
	var remoteType byte
	var seq int32
	var status int16

	buff := bytes.NewBuffer(raw)

	if err := binary.Read(buff, binary.BigEndian, &remoteVersion); err != nil {
		return CloseMessage{}, err
	}
	if remoteVersion != misc.ProtocolVersion {
		return CloseMessage{}, fmt.Errorf("Incorrect version (%v, should be %v)", remoteVersion, misc.ProtocolVersion)
	}

	if err := binary.Read(buff, binary.BigEndian, &remoteType); err != nil {
		return CloseMessage{}, err
	}
	if remoteType != byte(CloseMsgType) {
		return CloseMessage{}, fmt.Errorf("Incorrect type (%v, should be %v)", remoteType, byte(CloseMsgType))
	}

	if err := binary.Read(buff, binary.BigEndian, &seq); err != nil {
		return CloseMessage{}, err
	}
	if err := binary.Read(buff, binary.BigEndian, &status); err != nil {
		return CloseMessage{}, err
	}

	return CloseMessage{
		seq,
		status,
	}, nil
}

// Version returns the protocol version
func (m *CloseMessage) Version() byte {
	return misc.ProtocolVersion
}

// Type returns message type
func (m *CloseMessage) Type() MessageType {
	return CloseMsgType
}

// FixedLength returns the length of the fixed part of the msg
func (m *CloseMessage) FixedLength() uint16 {
	return 8
}

// Bytes returns the Close message in raw byte array
func (m *CloseMessage) Bytes() ([]byte, error) {
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
	if err := binary.Write(buff, binary.BigEndian, int16(m.status)); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
