package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"sip-monitor/src/pkg/siprocket"

	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/env"
	"sip-monitor/src/pkg/hep"
	"sip-monitor/src/pkg/parser"

	"github.com/sirupsen/logrus"
)

func HepServerListener() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: env.Conf.UDPListenPort})
	if err != nil {
		logrus.WithError(err).Error("HepServerListener Udp Service listen report udp fail")
	}

	defer conn.Close()
	logrus.Info("HepServerListener")

	var data = make([]byte, env.Conf.MaxPacketLength)
	var raw []byte
	for {
		err = conn.SetDeadline(time.Now().Add(time.Duration(env.Conf.MaxReadTimeoutSeconds) * time.Second))
		if err != nil {
			return
		}
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else {
				logrus.WithFields(logrus.Fields{
					"remote_addr": remoteAddr.IP.String(),
				}).WithError(err).Error("read udp error")
			}
		}

		if n < entity.MinRawPacketLength {
			logrus.WithFields(logrus.Fields{
				"setting_length":  entity.MinRawPacketLength,
				"received_length": n,
				"remote_addr":     remoteAddr.IP.String(),
			}).Warn("HepServerListener less then MinRawPacketLength")
			continue
		}

		raw = make([]byte, n)

		copy(raw, data[:n])

		go ParseSIPMsg(raw, remoteAddr.IP)
	}
}

func ParseSIPMsg(b []byte, ip net.IP) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err != nil {
			logrus.WithError(err.(error)).Error("parse save err")
		}
	}()

	s, errType, errMsg := Format(b)
	if errType != "" {
		logrus.WithFields(logrus.Fields{
			"raw_length": len(b),
			"from":       ip.String(),
			"err_type":   errType,
			"err_msg":    errMsg,
		}).Warn("format msg error")
		return
	}

	output := siprocket.Parse([]byte(*s.Raw))

	s.NodeIP = ip.String()
	s.TimestampMicroWithDate = s.CreateAt.Add(time.Microsecond * time.Duration(s.TimestampMicro)).UnixMicro()
	s.ViaNum = len(output.Via)

	model.SaveToDBQueue <- *s
}

func Format(p []byte) (s *entity.SIP, errorType string, errMsg string) {
	hepMsg, err := hep.NewHepMsg(p)

	if err != nil {
		logrus.WithError(err).Error("NewHepMsg error")
		return nil, "hep_parse_error", ""
	}

	if len(hepMsg.Body) <= 0 {
		return nil, "hep_body_is_empty", ""
	}

	if len(hepMsg.Body) < env.Conf.MinPacketLength {
		return nil, "hep_body_is_too_small", ""
	}

	sip := parser.Parser{SIP: entity.SIP{}}

	bodyS := string(hepMsg.Body)

	sip.SIP.Raw = &bodyS
	sip.ParseCseq()

	sip.TimestampMicro = hepMsg.TimestampMicro

	if sip.CSeqMethod == "" {
		return nil, "cseq_is_empty", ""
	}

	if strings.Contains(env.Conf.DiscardMethods, sip.CSeqMethod) {
		return nil, "method_discarded", sip.CSeqMethod
	}

	sip.ParseCallID()

	if sip.CallID == "" {
		return nil, "callid_is_empty", ""
	}

	sip.ParseFirstLine()

	if sip.Title == "" {
		return nil, "title_is_empty", ""
	}

	if sip.RequestURL != "" {
		sip.ParseRequestURL()
	}

	sip.ParseFrom()
	sip.ParseTo()
	sip.ParseUserAgent()
	sip.CreateAt = time.Unix(int64(hepMsg.Timestamp), 0)

	if env.Conf.HeaderFSCallIDName != "" {
		sip.ParseFSCallID(env.Conf.HeaderFSCallIDName)
	}

	if env.Conf.HeaderUIDName != "" {
		sip.ParseUID(env.Conf.HeaderUIDName)
	}

	sip.Protocol = int(hepMsg.IPProtocolID)

	sip.SrcAddr = fmt.Sprintf("%s_%d", hepMsg.IP4SourceAddress, hepMsg.SourcePort)
	sip.SrcPort = int(hepMsg.SourcePort)
	sip.SrcHost = hepMsg.IP4SourceAddress

	sip.DstAddr = fmt.Sprintf("%s_%d", hepMsg.IP4DestinationAddress, hepMsg.DestinationPort)
	sip.DstHost = hepMsg.IP4DestinationAddress
	sip.DstPort = int(hepMsg.DestinationPort)

	sip.NodeID = strconv.Itoa(int(hepMsg.CaptureAgentID))

	return &sip.SIP, "", ""
}

func BytesToInt(bys []byte) int {
	var data int64
	_ = binary.Read(bytes.NewBuffer(bys), binary.BigEndian, &data)
	return int(data)
}
