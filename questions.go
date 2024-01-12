// Question section format
//
// The question section is used to carry the "question" in most queries,
// i.e., the parameters that define what is being asked.  The section
// contains QDCOUNT (usually 1) entries, each of the following format:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                     QNAME                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QTYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QCLASS                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
// where:
//
// QNAME           a domain name represented as a sequence of labels, where
//                 each label consists of a length octet followed by that
//                 number of octets.  The domain name terminates with the
//                 zero length octet for the null label of the root.  Note
//                 that this field may be an odd number of octets; no
//                 padding is used.
//
// QTYPE           a two octet code which specifies the type of the query.
//                 The values for this field include all codes valid for a
//                 TYPE field, together with some more general codes which
//                 can match more than one type of RR.
//
// QCLASS          a two octet code that specifies the class of the query.
//                 For example, the QCLASS field is IN for the Internet.

package main

import (
	"bytes"
	"encoding/binary"
	"strings"

	bytes2 "github.com/thehivecorporation/bytes"
)

type Class uint16

const (
	_ Class = iota
	IN
	CS
	CH
	HS
)

type Type uint16

const (
	_ Type = iota
	A
	NS
	MD
	MF
	CNAME
	SOA
	MB
	MG
	MR
	NULL
	WKS
	PTR
	HINFO
	MINFO
	MX
	TXT
)

type Quesion struct {
	QName  string
	QType  Type
	QClass Class
}

func ReadQuesion(buf []byte) Quesion {
	start := 0
	var nameParts []string

	for len := buf[start]; len != 0; len = buf[start] {
		start++
		nameParts = append(nameParts, string(buf[start:start+int(len)]))
		start += int(len)
	}
	questionName := strings.Join(nameParts, ".")
	start++

	questionType := binary.BigEndian.Uint16(buf[start : start+2])
	questionClass := binary.BigEndian.Uint16(buf[start+2 : start+4])

	q := Quesion{
		QName:  questionName,
		QType:  Type(questionType),
		QClass: Class(questionClass),
	}

	return q
}

func (q Quesion) Encode() []byte {
	domain := q.QName
	parts := strings.Split(domain, ".")

	var buf bytes.Buffer

	for _, label := range parts {
		if len(label) > 0 {
			buf.WriteByte(byte(len(label)))
			buf.WriteString(label)
		}
	}

	buf.WriteByte(0x00)
	buf.Write(bytes2.Uint16ToBytes(uint16(q.QType)))
	buf.Write(bytes2.Uint16ToBytes(uint16(q.QClass)))

	return buf.Bytes()
}
