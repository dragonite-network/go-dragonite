package cryptor

// PacketCryptor encrypts and decrypts network packets
type PacketCryptor interface {
	Encrypt(plain []byte) []byte
	Decrypt(cipher []byte) []byte
	GetMaxAdditionalBytesLength() int
}
