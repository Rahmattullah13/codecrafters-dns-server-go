package main
import (
	"encoding/binary"
	"fmt"
)
type DNSHeader struct {
	id               uint16 //16 bits - the random id of query/reply
	isReply          bool   //1 bit  - true if reply false if query
	opCode           uint8  //4 bits - the type of query
	authoritative    bool   //1 bit - if reeplier owns domain
	truncated        bool   //1 bit - if message is larger than 512 bytes
	recursionDesired bool   //1 bit
	recursionAvail   bool   //1 bit
	reserved         uint8  //3 bits
	responseCode     uint8  //4 bits - kinda like a status code
	questionCount    uint16 //16 bits
	answerCount      uint16 //16 bits
	authCount        uint16 //16 bits
	additonalCount   uint16 //16 bits

}
func headerFromBytes(headerBytes []byte) DNSHeader {
	var id, flags, qCount, ansCount, auCount, adCount uint16
	id = binary.BigEndian.Uint16(headerBytes[0:2])
	flags = binary.BigEndian.Uint16(headerBytes[2:4])
	qCount = binary.BigEndian.Uint16(headerBytes[4:6])
	ansCount = binary.BigEndian.Uint16(headerBytes[6:8])
	auCount = binary.BigEndian.Uint16(headerBytes[8:10])
	adCount = binary.BigEndian.Uint16(headerBytes[10:12])
	return DNSHeader{
		id:               id,
		isReply:          (flags & 0x01) != 0,
		opCode:           uint8((flags & 0x1e) >> 1),
		authoritative:    ((flags & 0x20) >> 5) != 0,
		truncated:        ((flags & 0x40) >> 6) != 0,
		recursionDesired: ((flags & 0x80) >> 7) != 0,
		recursionAvail:   ((flags & 0x100) >> 8) != 0,
		reserved:         uint8((flags & 0xe00) >> 9),
		responseCode:     uint8((flags & 0xf000) >> 12),
		questionCount:    qCount,
		answerCount:      ansCount,
		authCount:        auCount,
		additonalCount:   adCount,

	}
}
type DNSMessage struct {
	header   DNSHeader
	contents string
}
func messageFromBytes(message []byte) DNSMessage {
	return DNSMessage{
		header:   headerFromBytes(message[0:12]),
		contents: string(message[12:]),

	}
}
func BoolToInt(b bool) uint16 {
	if b {
		return 1
	} else {
		return 0

	}
}
func (header DNSHeader) _flagsAsInt() uint16 {
	var flags uint16 = 0
	flags += uint16(header.responseCode)
	flags += uint16(header.reserved) << 4
	flags += BoolToInt(header.recursionAvail) << 7
	flags += BoolToInt(header.recursionDesired) << 8
	flags += BoolToInt(header.truncated) << 9
	flags += BoolToInt(header.authoritative) << 10
	flags += uint16(header.opCode) << 11
	flags += BoolToInt(header.isReply) << 15
	fmt.Println(flags)

	return flags
}
func (header DNSHeader) String() string {
	return fmt.Sprintf(
		"Id: %d\nisReply: %d\n opcode: %d\n authoritative: %d\n truncated: %d\n recursionD: %d\n recursionA: %d\n reserved: %d\n responseCode: %d qCount: %d\n ansCount: %d\n authCount: %d\n addCount: %d\n",
		header.id,
		BoolToInt(header.isReply),
		header.opCode,
		BoolToInt(header.authoritative),
		BoolToInt(header.truncated),
		BoolToInt(header.recursionDesired),
		BoolToInt(header.recursionAvail),
		header.reserved,
		header.responseCode,
		header.questionCount,
		header.answerCount,
		header.authCount,
		header.additonalCount,

	)
}
func (header DNSHeader) ToByteArray() []byte {
	headerBytes := make([]byte, 12)
	flags := header._flagsAsInt()
	binary.BigEndian.PutUint16(headerBytes[:2], header.id)
	binary.BigEndian.PutUint16(headerBytes[2:4], flags)
	binary.BigEndian.PutUint16(headerBytes[4:6], header.questionCount)
	binary.BigEndian.PutUint16(headerBytes[6:8], header.answerCount)
	binary.BigEndian.PutUint16(headerBytes[8:10], header.authCount)
	binary.BigEndian.PutUint16(headerBytes[10:], header.additonalCount)

	return headerBytes
}
func (message DNSMessage) ToByteArray() []byte {
	headerBytes := message.header.ToByteArray()
	contentBytes := []byte{}

	return append(headerBytes, contentBytes...)
}