// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// taken from http://golang.org/src/pkg/net/mockicmp_test.go

package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	icmpv4EchoRequest     = 8
	icmpv4EchoReply       = 0
	icmpv4DestUnreachable = 3
	icmpv6EchoRequest     = 128
	icmpv6EchoReply       = 129
)

func (m *icmpMessage) DestUnreachableMsg() string {
	if m.Type != icmpv4DestUnreachable {
		return ""
	}

	switch m.Code {
	case 0:
		return "destination network unreachable"
	case 1:
		return "destination host unreachable"
	case 2:
		return "destination protocol unreachable"
	case 3:
		return "destination port unreachable"
	case 4:
		return "fragmentation required and DF flag set"
	case 5:
		return "source route failed"
	case 6:
		return "destination network unknown"
	case 7:
		return "destination host unknown"
	case 8:
		return "source host isolated"
	case 9:
		return "network administratively prohibited"
	case 10:
		return "host administratively prohibited"
	case 11:
		return "network unreachable for TOS"
	case 12:
		return "host unreachable for TOS"
	case 13:
		return "communication administratively prohibited"
	case 14:
		return "host Precedence Violation"
	case 15:
		return "precedence cutoff in effect"
	}
	return "destination unreachable: unknown error"
}

// icmpMessage represents an ICMP message.
type icmpMessage struct {
	Type     int             // type
	Code     int             // code
	Checksum int             // checksum
	Body     icmpMessageBody // body
}

// icmpMessageBody represents an ICMP message body.
type icmpMessageBody interface {
	Len() int
	Marshal() ([]byte, error)
}

// Marshal returns the binary enconding of the ICMP echo request or
// reply message m.
func (m *icmpMessage) Marshal() ([]byte, error) {
	b := []byte{byte(m.Type), byte(m.Code), 0, 0}
	if m.Body != nil && m.Body.Len() != 0 {
		mb, err := m.Body.Marshal()
		if err != nil {
			return nil, err
		}
		b = append(b, mb...)
	}
	switch m.Type {
	case icmpv6EchoRequest, icmpv6EchoReply:
		return b, nil
	}
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	b[2] ^= byte(^s)
	b[3] ^= byte(^s >> 8)
	return b, nil
}

// parseICMPMessage parses b as an ICMP message.
func parseICMPMessage(b []byte) (*icmpMessage, error) {
	msglen := len(b)
	if msglen < 4 {
		return nil, errors.New("message too short")
	}
	m := &icmpMessage{Type: int(b[0]), Code: int(b[1]), Checksum: int(b[2])<<8 | int(b[3])}
	if msglen > 4 {
		var err error
		switch m.Type {
		case icmpv4EchoRequest, icmpv4EchoReply, icmpv6EchoRequest, icmpv6EchoReply:
			m.Body, err = parseICMPEcho(b[4:])
			if err != nil {
				return nil, err
			}
		}
	}
	return m, nil
}

// imcpEcho represenets an ICMP echo request or reply message body.
type icmpEcho struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

func (p *icmpEcho) Len() int {
	if p == nil {
		return 0
	}
	return 4 + len(p.Data)
}

// Marshal returns the binary enconding of the ICMP echo request or
// reply message body p.
func (p *icmpEcho) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(p.Data))
	b[0], b[1] = byte(p.ID>>8), byte(p.ID)
	b[2], b[3] = byte(p.Seq>>8), byte(p.Seq)
	copy(b[4:], p.Data)
	return b, nil
}

// parseICMPEcho parses b as an ICMP echo request or reply message
// body.
func parseICMPEcho(b []byte) (*icmpEcho, error) {
	bodylen := len(b)
	p := &icmpEcho{ID: int(b[0])<<8 | int(b[1]), Seq: int(b[2])<<8 | int(b[3])}
	if bodylen > 4 {
		p.Data = make([]byte, bodylen-4)
		copy(p.Data, b[4:])
	}
	return p, nil
}

func (p *icmpEcho) Decode() (time.Time, error) {
	t := time.Time{}
	err := t.UnmarshalBinary(p.Data[:15])
	if err != nil {
		fmt.Printf("icmpEcho.UnmarshalBinary problem: %v\n", err)
	}

	return t, err
}
