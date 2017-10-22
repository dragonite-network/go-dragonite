package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dragonite-network/go-dragonite/sdk/misc"
)

// AckMessage is a message signaling ACK
type AckMessage struct {
	seqList     []int32
	consumedSeq int32
}

// NewAckMessage constructs an ACK message from sequence list and consumed sequence count
func NewAckMessage(seqList []int32, consumedSeq int32) AckMessage {
	return AckMessage{
		seqList,
		consumedSeq,
	}
}

// NewAckMessageFromBytes constructs an ACK message from a raw msg
func NewAckMessageFromBytes(raw []byte) (AckMessage, error) {
	var remoteVersion byte
	var remoteType byte
	var consumedSeq int32
	var seqCount uint16
	var seqList []int32

	buff := bytes.NewBuffer(raw)

	if err := binary.Read(buff, binary.BigEndian, &remoteVersion); err != nil {
		return AckMessage{}, err
	}
	if remoteVersion != misc.ProtocolVersion {
		return AckMessage{}, fmt.Errorf("Incorrect version (%v, should be %v)", remoteVersion, misc.ProtocolVersion)
	}

	if err := binary.Read(buff, binary.BigEndian, &remoteType); err != nil {
		return AckMessage{}, err
	}
	if remoteType != byte(AckMsgType) {
		return AckMessage{}, fmt.Errorf("Incorrect type (%v, should be %v)", remoteType, byte(AckMsgType))
	}

	if err := binary.Read(buff, binary.BigEndian, &consumedSeq); err != nil {
		return AckMessage{}, err
	}
	if err := binary.Read(buff, binary.BigEndian, &seqCount); err != nil {
		return AckMessage{}, err
	}

	seqList = make([]int32, seqCount)
	for i := uint16(0); i < seqCount; i++ {
		var seq int32
		if err := binary.Read(buff, binary.BigEndian, &seq); err != nil {
			return AckMessage{}, err
		}
		seqList[i] = seq
	}

	return AckMessage{
		seqList,
		consumedSeq,
	}, nil
}

// Version returns the protocol version
func (m *AckMessage) Version() byte {
	return misc.ProtocolVersion
}

// Type returns message type
func (m *AckMessage) Type() MessageType {
	return AckMsgType
}

// FixedLength returns the length of the fixed part of the msg
func (m *AckMessage) FixedLength() uint16 {
	return 8
}

// Bytes returns the ACK message in raw byte array
func (m *AckMessage) Bytes() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	if err := binary.Write(buff, binary.BigEndian, byte(m.Version())); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, byte(m.Type())); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, int32(m.consumedSeq)); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, uint16(len(m.seqList))); err != nil {
		return nil, err
	}
	for _, seq := range m.seqList {
		if err := binary.Write(buff, binary.BigEndian, int32(seq)); err != nil {
			return nil, err
		}
	}
	return buff.Bytes(), nil
}
