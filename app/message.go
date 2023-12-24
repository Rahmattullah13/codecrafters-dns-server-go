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
		isReply:          (flags & 0x8000 >> 15) != 0,
		opCode:           uint8((flags & 0x7800) >> 11),
		authoritative:    ((flags & 0x400) >> 10) != 0,
		truncated:        ((flags & 0x200) >> 9) != 0,
		recursionDesired: ((flags & 0x100) >> 8) != 0,
		recursionAvail:   ((flags & 0x080) >> 7) != 0,
		reserved:         uint8((flags & 0x070) >> 4),
		responseCode:     uint8((flags & 0x00f)),
		questionCount:    qCount,
		answerCount:      ansCount,
		authCount:        auCount,
		additonalCount:   adCount,
	}
}

type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

type DNSRecord struct {
	Question   DNSQuestion
	TimeToLive uint32
	Length     uint16
	Data       []byte
}

func answerToQuestion(question DNSQuestion, data []byte, length uint16, ttl uint32) DNSRecord {
	return DNSRecord{
		Question:   question,
		TimeToLive: ttl,
		Length:     length,
		Data:       data,
	}
}

type DNSMessage struct {
	Header    DNSHeader
	Questions []DNSQuestion
	Answers   []DNSRecord
}

func messageFromBytes(message []byte) DNSMessage {
	header := headerFromBytes(message[0:12])
	questions := make([]DNSQuestion, header.questionCount)
	answers := make([]DNSRecord, header.answerCount)

	messagePointer := 12

	println("message:-------")

	for _, b := range message {
		println(b)
	}

	for questionNum := 0; questionNum < int(header.questionCount); questionNum++ {
		namePointer := messagePointer
		for message[namePointer] != 0x00 {
			namePointer++
		}

		questions[questionNum].Name = make([]byte, namePointer-messagePointer+1)
		copy(questions[questionNum].Name, message[messagePointer:namePointer+1])
		messagePointer = namePointer + 2

		questions[questionNum].Type = binary.BigEndian.Uint16(message[messagePointer : messagePointer+2])
		questions[questionNum].Class = binary.BigEndian.Uint16(message[messagePointer+2 : messagePointer+4])

		messagePointer += 4

		println("Name:")

		for _, b := range questions[questionNum].Name {
			println(b)
		}
	}

	return DNSMessage{
		Header:    header,
		Questions: questions,
		Answers:   answers,
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
		"Id: %d\nisReply: %d\n opcode: %d\n authoritative: %d\n truncated: %d\n recursionD: %d\n recursionA: %d\n reserved: %d\n responseCode: %d\n qCount: %d\n ansCount: %d\n authCount: %d\n addCount: %d\n",
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

func (question DNSQuestion) ToByteArray() []byte {
	questionBytes := make([]byte, 4)

	binary.BigEndian.PutUint16(questionBytes[:2], question.Type)
	binary.BigEndian.PutUint16(questionBytes[2:4], question.Class)

	return append(question.Name, questionBytes...)
}

func (record DNSRecord) ToByteArray() []byte {
	recordBytes := make([]byte, 6)

	binary.BigEndian.PutUint32(recordBytes[:4], record.TimeToLive)
	binary.BigEndian.PutUint16(recordBytes[4:6], record.Length)

	recordBytes = append(recordBytes, record.Data...)

	return append(record.Question.ToByteArray(), recordBytes...)
}

func (message DNSMessage) ToByteArray() []byte {
	headerBytes := message.Header.ToByteArray()
	contentBytes := []byte{}

	for _, question := range message.Questions {
		contentBytes = append(contentBytes, question.ToByteArray()...)
	}

	for _, answer := range message.Answers {
		contentBytes = append(contentBytes, answer.ToByteArray()...)
	}

	return append(headerBytes, contentBytes...)
}
