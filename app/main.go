package main

import (
	"fmt"
	"net"
	"strings"
)

// func respondToMessage(message DNSMessage) DNSMessage {
// 	var respCode uint8 = 0
// 	if message.Header.opCode != 0 {
// 		respCode = 4
// 	}

// 	var responseHeader DNSHeader = DNSHeader{
// 		id:               message.Header.id,
// 		isReply:          true,
// 		opCode:           message.Header.opCode,
// 		authoritative:    false,
// 		truncated:        false,
// 		recursionDesired: message.Header.recursionDesired,
// 		recursionAvail:   false,
// 		reserved:         0,
// 		responseCode:     respCode,
// 		questionCount:    1,
// 		answerCount:      1,
// 		authCount:        0,
// 		additonalCount:   0,
// 	}

// 	question := DNSQuestion{
// 		Name:  domaintoName("codecrafters.io"),
// 		Type:  1,
// 		Class: 1,
// 	}

// 	var response DNSMessage = DNSMessage{
// 		Header:    responseHeader,
// 		Questions: []DNSQuestion{question},
// 		Answers:   []DNSRecord{answerToQuestion(question, []byte{0x08, 0x08, 0x08, 0x08}, 4, 60)},
// 	}

// 	return response
// }

func nameToDomain(name []byte) string {
	var domainParts []string
	for i := 0; i < len(name); {
		labelLen := int(name[i])
		if labelLen == 0 {
			break
		}
		i++
		label := name[i : i+labelLen]
		domainParts = append(domainParts, string(label))
		i += labelLen
	}
	return strings.Join(domainParts, ".")
}


func respondToMessage(message DNSMessage) DNSMessage {
	var respCode uint8 = 0
	if message.Header.opCode != 0 {
		respCode = 4
	}

	// Check if there is at least one question in the received message
	if len(message.Questions) == 0 {
		// No question found, handle this case according to your requirements
		// For now, set a default domain
		return DNSMessage{}
	}

	// Extract the domain from the first question
	domain := nameToDomain(message.Questions[0].Name)

	var responseHeader DNSHeader = DNSHeader{
		id:               1234,
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

	question := DNSQuestion{
		Name:  domaintoName(domain),
		Type:  1,
		Class: 1,
	}

	var response DNSMessage = DNSMessage{
		Header:    responseHeader,
		Questions: []DNSQuestion{question},
		Answers:   []DNSRecord{answerToQuestion(question, []byte{0x08, 0x08, 0x08, 0x08}, 4, 60)},
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
