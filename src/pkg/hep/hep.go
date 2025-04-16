/**
* Homer Encapsulation Protocol v3
* Courtesy of Weave Communications, Inc (http://getweave.com/) under the ISC license (https://en.wikipedia.org/wiki/ISC_license)
* https://www.voztovoice.org/sites/default/files/hep3_rev12.pdf
**/

package hep

import (
	"encoding/binary"
	"errors"
	"net"
)

/*************************************
 Constants
*************************************/

// HEP ID
const (
	HEPID1 = 0x011002
	HEPID2 = 0x021002
	HEPID3 = 0x48455033
)

// Generic Chunk Types
const (
	_ = iota // Don't want to assign zero here, but want to implicitly repeat this expression after...
	IPProtocolFamily
	IPProtocolID
	IP4SourceAddress
	IP4DestinationAddress
	IP6SourceAddress
	IP6DestinationAddress
	SourcePort
	DestinationPort
	Timestamp
	TimestampMicro
	ProtocolType // Maps to Protocol Types below
	CaptureAgentID
	KeepAliveTimer
	AuthenticationKey
	PacketPayload
	CompressedPayload
	InternalC
)

var protocolFamilies []string
var vendors []string
var protocolTypes []string

const (
	// 0x01 = SIP
	// 0x02 = XMPP
	// 0x03 = SDP
	// 0x04 = RTP
	// 0x05 = RTCP
	ProtocolTypeSIP  = 0x01
	ProtocolTypeXMPP = 0x02
	ProtocolTypeSDP  = 0x03
	ProtocolTypeRTP  = 0x04
	ProtocolTypeRTCP = 0x05
)

func init() {

	// Protocol Family Types - HEP3 Spec does not list these values out. Took IPv4 from an example.
	protocolFamilies = []string{
		"?",
		"?",
		"IPv4"}

	// Initialize vendors
	vendors = []string{
		"None",
		"FreeSWITCH",
		"Kamailio",
		"OpenSIPS",
		"Asterisk",
		"Homer",
		"SipXecs",
	}

	// Initialize protocol types
	protocolTypes = []string{
		"Reserved",
		"SIP",
		"XMPP",
		"SDP",
		"RTP",
		"RTCP",
		"MGCP",
		"MEGACO",
		"M2UA",
		"M3UA",
		"IAX",
		"H322",
		"H321",
	}
}

// HepMsg represents a parsed HEP packet
type HepMsg struct {
	IPProtocolFamily      byte
	IPProtocolID          byte
	IP4SourceAddress      string
	IP4DestinationAddress string
	IP6SourceAddress      string
	IP6DestinationAddress string
	SourcePort            uint16
	DestinationPort       uint16
	Timestamp             uint32
	TimestampMicro        uint32
	ProtocolType          byte
	CaptureAgentID        uint16
	KeepAliveTimer        uint16
	AuthenticateKey       string
	InternalCorrelationID string // 内部关联ID
	Body                  []byte
}

// NewHepMsg returns a parsed message object. Takes a byte slice.
func NewHepMsg(packet []byte) (*HepMsg, error) {
	newHepMsg := &HepMsg{}
	err := newHepMsg.parse(packet)
	if err != nil {
		return nil, err
	}
	return newHepMsg, nil
}

func NewMockHepMsgBySIPMsg(sipBody []byte) *HepMsg {
	return &HepMsg{
		IPProtocolFamily:      0x01,
		IPProtocolID:          0x01,
		IP4SourceAddress:      "127.0.0.1",
		IP4DestinationAddress: "127.0.0.1",
		SourcePort:            5060,
		DestinationPort:       5060,
		Timestamp:             0,
		TimestampMicro:        0,
		ProtocolType:          0x01,
		CaptureAgentID:        0x01,
		KeepAliveTimer:        0x01,
		AuthenticateKey:       "1234567890",
		Body:                  sipBody,
	}
}

func (hepMsg *HepMsg) parse(udpPacket []byte) error {
	switch udpPacket[0] {
	case 0x01:
		return hepMsg.parseHep1(udpPacket)
	case 0x02:
		return hepMsg.parseHep2(udpPacket)
	case 0x48:
		return hepMsg.parseHep3(udpPacket)
	default:
		err := errors.New("Not a valid HEP packet - HEP ID does not match spec")
		return err
	}
}
func (hepMsg *HepMsg) parseHep1(udpPacket []byte) error {
	//var err error
	if len(udpPacket) < 21 {
		return errors.New("Found HEP ID for HEP v1, but length of packet is too short to be HEP1 or is NAT keepalive")
	}
	packetLength := len(udpPacket)
	hepMsg.SourcePort = binary.BigEndian.Uint16(udpPacket[4:6])
	hepMsg.DestinationPort = binary.BigEndian.Uint16(udpPacket[6:8])
	hepMsg.IP4SourceAddress = net.IP(udpPacket[8:12]).String()
	hepMsg.IP4DestinationAddress = net.IP(udpPacket[12:16]).String()
	hepMsg.Body = udpPacket[16:]
	if len(udpPacket[16:packetLength-4]) > 1 {
		//SIP消息：udpPacket[16:packetLength]
	} else {

	}

	return nil
}

func (hepMsg *HepMsg) parseHep2(udpPacket []byte) error {
	//var err error
	if len(udpPacket) < 31 {
		return errors.New("Found HEP ID for HEP v2, but length of packet is too short to be HEP2 or is NAT keepalive")
	}
	packetLength := len(udpPacket)
	hepMsg.SourcePort = binary.BigEndian.Uint16(udpPacket[4:6])
	hepMsg.DestinationPort = binary.BigEndian.Uint16(udpPacket[6:8])
	hepMsg.IP4SourceAddress = net.IP(udpPacket[8:12]).String()
	hepMsg.IP4DestinationAddress = net.IP(udpPacket[12:16]).String()
	hepMsg.Timestamp = binary.LittleEndian.Uint32(udpPacket[16:20])
	hepMsg.TimestampMicro = binary.LittleEndian.Uint32(udpPacket[20:24])
	hepMsg.CaptureAgentID = binary.BigEndian.Uint16(udpPacket[24:26])
	hepMsg.Body = udpPacket[28:]
	if len(udpPacket[28:packetLength-4]) > 1 {
		//SIP消息：udpPacket[16:packetLength]
	} else {

	}

	return nil
}

func (hepMsg *HepMsg) parseHep3(udpPacket []byte) error {
	length := binary.BigEndian.Uint16(udpPacket[4:6])
	currentByte := uint16(6)

	for currentByte < length {
		hepChunk := udpPacket[currentByte:]
		//chunkVendorId := binary.BigEndian.Uint16(hepChunk[:2])
		chunkType := binary.BigEndian.Uint16(hepChunk[2:4])
		chunkLength := binary.BigEndian.Uint16(hepChunk[4:6])

		//todo::实际运行过程中，chunkLength会超过len(hepChunk)，导致下面报错
		if int(chunkLength) > len(hepChunk) {
			continue
		}

		chunkBody := hepChunk[6:chunkLength]

		switch chunkType {
		case IPProtocolFamily:
			hepMsg.IPProtocolFamily = chunkBody[0]
		case IPProtocolID:
			hepMsg.IPProtocolID = chunkBody[0]
		case IP4SourceAddress:
			hepMsg.IP4SourceAddress = net.IP(chunkBody).String()
		case IP4DestinationAddress:
			hepMsg.IP4DestinationAddress = net.IP(chunkBody).String()
		case IP6SourceAddress:
			hepMsg.IP6SourceAddress = net.IP(chunkBody).String()
		case IP6DestinationAddress:
			hepMsg.IP4DestinationAddress = net.IP(chunkBody).String()
		case SourcePort:
			hepMsg.SourcePort = binary.BigEndian.Uint16(chunkBody)
		case DestinationPort:
			hepMsg.DestinationPort = binary.BigEndian.Uint16(chunkBody)
		case Timestamp:
			hepMsg.Timestamp = binary.BigEndian.Uint32(chunkBody)
		case TimestampMicro:
			hepMsg.TimestampMicro = binary.BigEndian.Uint32(chunkBody)
		case ProtocolType:
			hepMsg.ProtocolType = chunkBody[0]
		case CaptureAgentID:
			hepMsg.CaptureAgentID = binary.BigEndian.Uint16(chunkBody)
		case KeepAliveTimer:
			hepMsg.KeepAliveTimer = binary.BigEndian.Uint16(chunkBody)
		case AuthenticationKey:
			hepMsg.AuthenticateKey = string(chunkBody)
		case PacketPayload:
			hepMsg.Body = append(hepMsg.Body, chunkBody...)
		case CompressedPayload:
		case InternalC:
			hepMsg.InternalCorrelationID = string(chunkBody)
		default:
		}
		currentByte += chunkLength
	}
	return nil
}
