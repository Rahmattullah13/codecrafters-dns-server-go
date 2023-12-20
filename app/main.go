package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
)

func createDNSHeader() []byte {
	packetID := uint16(1234)
	qr := byte(1) << 1
	opcode := byte(0)
	aa := byte(0)
	tc := byte(0)
	rd := byte(0)
	ra := byte(0)
	z := byte(0)
	rcode := byte(0)
	qdcount := uint16(0)
	ancount := uint16(0)
	nscount := uint16(0)
	arcount := uint16(0)

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

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

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

		// create a dns response with the specified header
		responseHeader := createDNSHeader()

		// combine the header and any additional response data
		response := append(responseHeader, []byte("Your additional response data here")...)

		// Create an empty response
		response = []byte{}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
