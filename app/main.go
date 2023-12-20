package main

import (
	"fmt"
	"net"
)

func respondToMessage(message DNSMessage) DNSMessage {
	var responseHeader DNSHeader = DNSHeader{
		id:               message.header.id,
		isReply:          true,
		opCode:           0,
		authoritative:    false,
		truncated:        false,
		recursionDesired: false,
		recursionAvail:   false,
		reserved:         0,
		responseCode:     0,
		questionCount:    0,
		answerCount:      0,
		authCount:        0,
		additonalCount:   0,
	}
	var response DNSMessage = DNSMessage{
		header:   responseHeader,
		contents: "",
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

		fmt.Println(messageBack.header.String())
		

		response := messageBack.ToByteArray()
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
