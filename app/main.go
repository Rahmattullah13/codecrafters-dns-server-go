package main

import (
	"fmt"
	"net"
)

func createDNSHeader() []byte {
	// Header fields
	packetID := uint16(1234)
	qr := byte(1) << 7 // 1 for response packet
	opcode := byte(0)
	aa := byte(0)
	tc := byte(0)
	rd := byte(0)
	ra := byte(0)
	z := byte(0)
	rcode := byte(0)
	qdcount := uint16(0)
	ancount := uint16(1) // Assuming one answer record
	nscount := uint16(0)
	arcount := uint16(0)

	// Pack the header fields into a byte slice
	header := []byte{
		byte(packetID >> 8), byte(packetID),
		qr | (opcode << 3) | (aa << 2) | (tc << 1) | rd,
		(ra << 7) | (z << 4) | rcode,
		byte(qdcount >> 8), byte(qdcount),
		byte(ancount >> 8), byte(ancount),
		byte(nscount >> 8), byte(nscount),
		byte(arcount >> 8), byte(arcount),
	}

	return header
}

func createAnswerRecord() []byte {
	// Sample answer record
	name := []byte{0xc0, 0x0c}            // Compressed name pointer to the domain name in the question
	qtype := []byte{0x00, 0x01}           // A record type
	qclass := []byte{0x00, 0x01}          // IN class
	ttl := []byte{0x00, 0x00, 0x00, 0x0a} // 10 seconds
	rdlength := []byte{0x00, 0x04}        // Length of the RData field (IPv4 address)

	// IPv4 address (replace with your actual IP address)
	rdata := net.ParseIP("127.0.0.1").To4()

	// Combine all components into the answer record
	answerRecord := append(name, qtype...)
	answerRecord = append(answerRecord, qclass...)
	answerRecord = append(answerRecord, ttl...)
	answerRecord = append(answerRecord, rdlength...)
	answerRecord = append(answerRecord, rdata...)

	return answerRecord
}

func main() {
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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create a DNS response with the specified header and answer record
		responseHeader := createDNSHeader()
		answerRecord := createAnswerRecord()
		response := append(responseHeader, answerRecord...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
