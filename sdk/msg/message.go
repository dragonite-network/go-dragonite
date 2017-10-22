package msg

// Message is a dragonite message
type Message interface {
	Version() byte
	Type() MessageType
	FixedLength() uint16
	Bytes() ([]byte, error)
}
