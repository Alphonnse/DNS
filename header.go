// The header contains the following fields:
//
//	                                1  1  1  1  1  1
//	  0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|                      ID                       |
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|                    QDCOUNT                    |
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|                    ANCOUNT                    |
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|                    NSCOUNT                    |
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|                    ARCOUNT                    |
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
// where:
//
// ID              A 16 bit identifier assigned by the program that
//
//	generates any kind of query.  This identifier is copied
//	the corresponding reply and can be used by the requester
//	to match up replies to outstanding queries.
//
// QR              A one bit field that specifies whether this message is a
//
//	query (0), or a response (1).
//
// OPCODE          A four bit field that specifies kind of query in this
//
//	message.  This value is set by the originator of a query
//	and copied into the response.  The values are:
//
//	0               a standard query (QUERY)
//
//	1               an inverse query (IQUERY)
//
//	2               a server status request (STATUS)
//
//	3-15            reserved for future use
//
// AA              Authoritative Answer - this bit is valid in responses,
//
//	and specifies that the responding name server is an
//	authority for the domain name in question section.
//
//	Note that the contents of the answer section may have
//	multiple owner names because of aliases.  The AA bit
package main

import "encoding/binary"

type Header struct {
	PacketID uint16
	QR       uint16
	OPCODE   uint16
	AA       uint16
	TC       uint16
	RD       uint16
	RA       uint16
	Z        uint16
	RCode    uint16
	QDCount  uint16
	ANCount  uint16
	NSCount  uint16
	ARCount  uint16
}

func ReadHeader(buf []byte) Header {
	h := Header{
		PacketID: uint16(buf[0])<<8 | uint16(buf[1]),
		QR:       1,
		OPCODE:   uint16(buf[2]<<1) >> 4,
		AA:       uint16(buf[3]<<5) >> 7,
		TC:       uint16(buf[3]<<6) >> 7,
		RD:       uint16(buf[3]<<7) >> 7,
		RA:       uint16(buf[3]) >> 7,
		Z:        uint16(buf[4]<<1) >> 5,
		QDCount:  uint16(buf[4])<<8 | uint16(buf[5]),
		ANCount:  uint16(buf[5])<<8 | uint16(buf[7]),
		NSCount:  uint16(buf[8])<<8 | uint16(buf[9]),
		ARCount:  uint16(buf[10])<<8 | uint16(buf[11]),
	}

	// I'll process only opcode 0 (a standard query (QUERY))
	if h.OPCODE == 0 {
		h.RCode = 0
	} else {
		h.RCode = 4
	}

	return h
}

func (h Header) Encode() []byte {
	dnsHeader := make([]byte, 12) // header is 12 bytes

	var flags uint16 = 0
	flags = h.QR<<15 | h.OPCODE<<11 | h.AA<<10 | h.TC<<9 | h.RD<<8 | h.RA<<7 | h.Z<<4 | h.RCode

	binary.BigEndian.PutUint16(dnsHeader[0:2], h.PacketID)
	binary.BigEndian.PutUint16(dnsHeader[2:4], flags)
	binary.BigEndian.PutUint16(dnsHeader[4:6], h.QDCount)
	binary.BigEndian.PutUint16(dnsHeader[6:8], h.ANCount)
	binary.BigEndian.PutUint16(dnsHeader[8:10], h.NSCount)
	binary.BigEndian.PutUint16(dnsHeader[10:12], h.ARCount)

	return dnsHeader
}
