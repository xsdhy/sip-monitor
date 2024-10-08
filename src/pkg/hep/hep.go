/**
* Homer Encapsulation Protocol v3
* Courtesy of Weave Communications, Inc (http://getweave.com/) under the ISC license (https://en.wikipedia.org/wiki/ISC_license)
**/

package hep

import (
	"encoding/binary"
	"errors"
	"net"
)

// HEP Versions
const (
	HEP1 byte = 0x01
	HEP2 byte = 0x02
	HEP3 byte = 0x48
)

// HEP3 Chunk Types
const (
	_ = iota
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

// HepMsg represents a parsed HEP packet
type HepMsg struct {
	Version byte   // HEP 协议版本，通常为 3
	Type    byte   // HEP 消息类型，例如 0 表示 HEP_DATA
	SubType uint16 // 消息子类型，根据 Type 不同而不同

	Timestamp      uint32 // Unix 时间戳（秒）
	TimestampMicro uint32 // Unix 时间戳的微秒部分

	CaptureAgentID uint16 // 捕获代理的唯一标识符

	IPProtocolFamily byte // 地址族：1 表示 IPv4，2 表示 IPv6
	IPProtocolID     byte // 上层协议标识符，例如 6 表示 TCP，17 表示 UDP

	IP4SourceAddress      string // 源 IPv4 地址
	IP4DestinationAddress string // 目标 IPv4 地址

	IP6SourceAddress      string // 源 IPv6 地址
	IP6DestinationAddress string // 目标 IPv6 地址

	SourcePort      uint16 // 源端口号
	DestinationPort uint16 // 目标端口号

	ProtocolType byte // 应用层协议类型，例如 SIP、RTP

	KeepAliveTimer  uint16 // 保持连接活跃的定时器值
	AuthenticateKey string // 认证密钥（固定长度）
	Body            []byte // 消息的有效负载
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

func (hepMsg *HepMsg) parse(udpPacket []byte) error {
	switch udpPacket[0] {
	case HEP1:
		return hepMsg.parseHep1(udpPacket)
	case HEP2:
		return hepMsg.parseHep2(udpPacket)
	case HEP3:
		return hepMsg.parseHep3(udpPacket)
	default:
		return errors.New("not a valid HEP packet - HEP ID does not match spec")
	}
}
func (hepMsg *HepMsg) parseHep1(udpPacket []byte) error {
	//var err error
	if len(udpPacket) < 21 {
		return errors.New("found HEP ID for HEP v1, but length of packet is too short to be HEP1 or is NAT keepalive")
	}
	packetLength := len(udpPacket)
	hepMsg.SourcePort = binary.BigEndian.Uint16(udpPacket[4:6])
	hepMsg.DestinationPort = binary.BigEndian.Uint16(udpPacket[6:8])
	hepMsg.IP4SourceAddress = net.IP(udpPacket[8:12]).String()
	hepMsg.IP4DestinationAddress = net.IP(udpPacket[12:16]).String()
	hepMsg.Body = udpPacket[16:]
	if len(udpPacket[16:packetLength-4]) > 1 {
		//hepMsg.SipMsg = siprocket.Parse(udpPacket[16:packetLength])
	} else {

	}

	return nil
}

func (hepMsg *HepMsg) parseHep2(udpPacket []byte) error {
	//var err error
	if len(udpPacket) < 31 {
		return errors.New("found HEP ID for HEP v2, but length of packet is too short to be HEP2 or is NAT keepalive")
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
		//hepMsg.SipMsg = siprocket.Parse(udpPacket[16:packetLength])
	} else {

	}

	return nil
}

func (hepMsg *HepMsg) parseHep3(udpPacket []byte) error {
	if len(udpPacket) < 6 {
		return errors.New("HEP3 packet too short to contain length field")
	}
	hepMsg.Version = 3

	length := binary.BigEndian.Uint16(udpPacket[4:6])
	currentByte := uint16(6)

	for currentByte < length {
		hepChunk := udpPacket[currentByte:]
		//chunkVendorId := binary.BigEndian.Uint16(hepChunk[:2])
		chunkType := binary.BigEndian.Uint16(hepChunk[2:4])
		chunkLength := binary.BigEndian.Uint16(hepChunk[4:6])

		//实际运行过程中，chunkLength会超过len(hepChunk)，所以需要有这个检查
		if int(chunkLength) > len(hepChunk) {
			return errors.New("HEP3 packet too short to contain length field")
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
			hepMsg.IP6DestinationAddress = net.IP(chunkBody).String()
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
		default:
		}
		currentByte += chunkLength
	}
	return nil
}
