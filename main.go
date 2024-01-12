// All communications inside of the domain protocol are carried in a single
// format called a message.  The top level format of message is divided
// into 5 sections (some of which are empty in certain cases) shown below:
//
//     +---------------------+
//     |        Header       |
//     +---------------------+
//     |       Question      | the question for the name server
//     +---------------------+
//     |        Answer       | RRs answering the question
//     +---------------------+
//     |      Authority      | RRs pointing toward an authority
//     +---------------------+
//     |      Additional     | RRs holding additional information
//     +---------------------+

package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)


const address = "127.0.0.1:2053"

var nameToIp = map[string][4]uint8 {
	"google.com": {12, 32, 31, 12},
	"gnu.org": {33, 3, 3, 33},
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal("Failed to resolve udp address", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to bind to address", err)
	}
	
	defer udpConn.Close()

	log.Printf("started server on %s", address)

	// 512 byte buffer (RFC)
	buf := make([]byte, 512)
	for {
		_, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Failed to receive data", err)
			break
		}
		
		header := ReadHeader(buf[:12])
		question := ReadQuesion(buf[12:])

		fmt.Println(question.QName) // у меня что то в question плывет

		answer := Answer {
			Name: question.QName,
			Type: A,
			Class: IN,
			TTL: 0,
			Rdlenght: net.IPv4len,
			Rdata: nameToIp[question.QName],
		}

		var response bytes.Buffer

		response.Write(header.Encode())
		response.Write(question.Encode())
		response.Write(answer.Encode())

		_, err = udpConn.WriteToUDP(response.Bytes(), source)
		if err != nil {
			log.Println("Failed to send response:", err)
		}
	}
}

