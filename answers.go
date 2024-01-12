// The answer, authority, and additional sections all share the same
// format: a variable number of resource records, where the number of
// records is specified in the corresponding count field in the header.
// Each resource record has the following format:
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                                               /
//     /                      NAME                     /
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     CLASS                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TTL                      |
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                   RDLENGTH                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//     /                     RDATA                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
// where:
//
// NAME            a domain name to which this resource record pertains.
//
// TYPE            two octets containing one of the RR type codes.  This
//                 field specifies the meaning of the data in the RDATA
//                 field.
//
// CLASS           two octets which specify the class of the data in the
//                 RDATA field.
//
// TTL             a 32 bit unsigned integer that specifies the time
//                 interval (in seconds) that the resource record may be
//                 cached before it should be discarded.  Zero values are
//                 interpreted to mean that the RR can only be used for the
//                 transaction in progress, and should not be cached.
//
// RDLENGTH        an unsigned 16 bit integer that specifies the length in
//                 octets of the RDATA field.
//
// RDATA           a variable length string of octets that describes the
//                 resource.  The format of this information varies
//                 according to the TYPE and CLASS of the resource record.
//                 For example, the if the TYPE is A and the CLASS is IN,
//                 the RDATA field is a 4 octet ARPA Internet address.

package main

import (
	"encoding/binary"
	"net"
	"strings"

	"github.com/thehivecorporation/bytes"
)

type Answer struct {
	Name string
	Type Type
	Class Class
	TTL uint32
	Rdlenght uint16
	Rdata [4]uint8
}

func (a Answer) Encode() []byte{
	var buff []byte

	domain := a.Name
	parts := strings.Split(domain, ".")

	for _, label := range parts {
		if len(label) > 0 {
			buff = append(buff, byte(len(label)))
			buff = append(buff, []byte(label)...)
		}
	}

	buff = append(buff, 0x00)

	buff = append(buff, bytes.Uint16ToBytes(uint16(a.Type))...)
	buff = append(buff, bytes.Uint16ToBytes(uint16(a.Class))...)

	time := make([]byte,4)
	binary.BigEndian.PutUint32(time, a.TTL)

	buff = append(buff, time...)
	buff = append(buff, bytes.Uint16ToBytes(a.Rdlenght)...)

	ipBytes, err := net.IPv4(a.Rdata[0], a.Rdata[1], a.Rdata[2], a.Rdata[3]).MarshalText()
	if err != nil {
		return nil
	}

	buff = append(buff, ipBytes...)

	return buff
}
