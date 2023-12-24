package main

import (
	"fmt"
	"net"
)

func respondToMessage(message DNSMessage) DNSMessage {
	var respCode uint8 = 0
	if message.Header.opCode != 0 {
		respCode = 4
	}

	var responseHeader DNSHeader = DNSHeader{
		id:               message.Header.id,
		isReply:          true,
		opCode:           message.Header.opCode,
		authoritative:    false,
		truncated:        false,
		recursionDesired: message.Header.recursionDesired,
		recursionAvail:   false,
		reserved:         0,
		responseCode:     respCode,
		questionCount:    1,
		answerCount:      1,
		authCount:        0,
		additonalCount:   0,
	}

	questions := make([]DNSQuestion, message.Header.questionCount)
	answers := make([]DNSRecord, message.Header.questionCount)
	for index, question := range message.Questions {
		questions[index].Name = question.Name
		questions[index].Type = 1
		questions[index].Class = 1
		answers[index] = answerToQuestion(questions[index], []byte{0x08, 0x08, 0x08, 0x08}, 4, 60)
	}

	var response DNSMessage = DNSMessage{
		Header:    responseHeader,
		Questions: questions,
		Answers:   answers,
	}

	return response
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, string(receivedData))
		receivedMessage := messageFromBytes(receivedData)
		messageBack := respondToMessage(receivedMessage)

		fmt.Println(receivedMessage.Header.String())

		response := messageBack.ToByteArray()
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
