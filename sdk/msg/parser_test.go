package msg

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/dragonite-network/go-dragonite/sdk/misc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDataMessage(t *testing.T) {
	seqBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(seqBuff, binary.BigEndian, int32(777)))
	dataLenBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(dataLenBuff, binary.BigEndian, uint16(4)))
	data := []byte{1, 2, 3, 4}

	raw := []byte{misc.ProtocolVersion, byte(DataMsgType)}
	raw = append(raw, seqBuff.Bytes()...)
	raw = append(raw, dataLenBuff.Bytes()...)
	raw = append(raw, data...)

	parsedMessage, err := ParseMessage(raw)
	require.NoError(t, err)
	assert.Equal(t, DataMsgType, parsedMessage.Type())
	assert.Equal(t, misc.ProtocolVersion, parsedMessage.Version())

	// data message specs
	dataMessage := parsedMessage.(*DataMessage)
	assert.Equal(t, data, dataMessage.data)
	assert.EqualValues(t, 777, dataMessage.seq)

	// roundtrip
	toRaw, err := dataMessage.Bytes()
	require.NoError(t, err)
	assert.Equal(t, raw, toRaw)
}

func TestParseAckMessage(t *testing.T) {
	seqBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(seqBuff, binary.BigEndian, int32(777)))
	consumedSeqBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(consumedSeqBuff, binary.BigEndian, int32(42)))
	seqCountBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(seqCountBuff, binary.BigEndian, uint16(4)))

	seqListBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(seqListBuff, binary.BigEndian, int32(5)))
	require.NoError(t, binary.Write(seqListBuff, binary.BigEndian, int32(6)))
	require.NoError(t, binary.Write(seqListBuff, binary.BigEndian, int32(7)))
	require.NoError(t, binary.Write(seqListBuff, binary.BigEndian, int32(8)))

	raw := []byte{misc.ProtocolVersion, byte(AckMsgType)}
	raw = append(raw, consumedSeqBuff.Bytes()...)
	raw = append(raw, seqCountBuff.Bytes()...)
	raw = append(raw, seqListBuff.Bytes()...)

	parsedMessage, err := ParseMessage(raw)
	require.NoError(t, err)
	assert.Equal(t, AckMsgType, parsedMessage.Type())
	assert.Equal(t, misc.ProtocolVersion, parsedMessage.Version())

	// ack message specs
	ackMessage := parsedMessage.(*AckMessage)
	assert.EqualValues(t, 42, ackMessage.consumedSeq)
	assert.EqualValues(t, []int32{5, 6, 7, 8}, ackMessage.seqList)

	// roundtrip
	toRaw, err := parsedMessage.Bytes()
	require.NoError(t, err)
	assert.Equal(t, raw, toRaw)
}

func TestParseCloseMessage(t *testing.T) {
	seqBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(seqBuff, binary.BigEndian, int32(777)))
	statusBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(statusBuff, binary.BigEndian, int16(42)))

	raw := []byte{misc.ProtocolVersion, byte(CloseMsgType)}
	raw = append(raw, seqBuff.Bytes()...)
	raw = append(raw, statusBuff.Bytes()...)

	parsedMessage, err := ParseMessage(raw)
	require.NoError(t, err)
	assert.Equal(t, CloseMsgType, parsedMessage.Type())
	assert.Equal(t, misc.ProtocolVersion, parsedMessage.Version())

	// close message specs
	closeMessage := parsedMessage.(*CloseMessage)
	assert.EqualValues(t, 777, closeMessage.seq)
	assert.EqualValues(t, 42, closeMessage.status)

	// roundtrip
	toRaw, err := parsedMessage.Bytes()
	require.NoError(t, err)
	assert.Equal(t, raw, toRaw)
}

func TestParseHeartbeatMessage(t *testing.T) {
	seqBuff := bytes.NewBuffer([]byte{})
	require.NoError(t, binary.Write(seqBuff, binary.BigEndian, int32(777)))

	raw := []byte{misc.ProtocolVersion, byte(HeartbeatMsgType)}
	raw = append(raw, seqBuff.Bytes()...)

	parsedMessage, err := ParseMessage(raw)
	require.NoError(t, err)
	assert.Equal(t, HeartbeatMsgType, parsedMessage.Type())
	assert.Equal(t, misc.ProtocolVersion, parsedMessage.Version())

	// heartbeat message specs
	heartBeatMessage := parsedMessage.(*HeartbeatMessage)
	assert.EqualValues(t, 777, heartBeatMessage.seq)

	// roundtrip
	toRaw, err := parsedMessage.Bytes()
	require.NoError(t, err)
	assert.Equal(t, raw, toRaw)
}
