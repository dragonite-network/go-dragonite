package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/dragonite-network/go-dragonite/sdk/misc"
)

// DataMessage is a message carrying a data payload
type DataMessage struct {
	seq  int32
	data []byte
}

// NewDataMessage constructs a data message from sequence and data
func NewDataMessage(seq int32, data []byte) DataMessage {
	return DataMessage{
		seq,
		data,
	}
}

// NewDataMessageFromBytes constructs a data message from raw msg bytes
func NewDataMessageFromBytes(raw []byte) (DataMessage, error) {
	var remoteVersion byte
	var remoteType byte
	var seq int32
	var dataLen uint16
	var data []byte

	buff := bytes.NewBuffer(raw)

	if err := binary.Read(buff, binary.BigEndian, &remoteVersion); err != nil {
		return DataMessage{}, err
	}
	if remoteVersion != misc.ProtocolVersion {
		return DataMessage{}, fmt.Errorf("Incorrect version (%v, should be %v)", remoteVersion, misc.ProtocolVersion)
	}
	if err := binary.Read(buff, binary.BigEndian, &remoteType); err != nil {
		return DataMessage{}, err
	}
	if remoteType != byte(DataMsgType) {
		return DataMessage{}, fmt.Errorf("Incorrect type (%v, should be %v)", remoteType, byte(DataMsgType))
	}
	if err := binary.Read(buff, binary.BigEndian, &seq); err != nil {
		return DataMessage{}, err
	}
	if err := binary.Read(buff, binary.BigEndian, &dataLen); err != nil {
		return DataMessage{}, err
	}

	data = make([]byte, dataLen)
	if _, err := io.ReadAtLeast(buff, data, int(dataLen)); err != nil {
		return DataMessage{}, err
	}

	return DataMessage{
		seq,
		data,
	}, nil
}

// Version returns the protocol version
func (m *DataMessage) Version() byte {
	return misc.ProtocolVersion
}

// Type returns message type
func (m *DataMessage) Type() MessageType {
	return DataMsgType
}

// FixedLength returns the length of the fixed part of the msg
func (m *DataMessage) FixedLength() uint16 {
	return 8
}

// Bytes returns the message in raw byte array
func (m *DataMessage) Bytes() ([]byte, error) {
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
	if err := binary.Write(buff, binary.BigEndian, uint16(len(m.data))); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, m.data); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
