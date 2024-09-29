package services

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/env"
	"sip-monitor/src/pkg/hep"
	"sip-monitor/src/pkg/parser"
)

func HepServerListener() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: env.Conf.UDPListenPort})
	if err != nil {
		slog.Error("HepServerListener Udp Service listen report udp fail", slog.String("reason", err.Error()))
	}

	defer conn.Close()
	slog.Info("HepServerListener")

	var data = make([]byte, env.Conf.MaxPacketLength)
	var raw []byte
	for {
		err = conn.SetDeadline(time.Now().Add(time.Duration(env.Conf.MaxReadTimeoutSeconds) * time.Second))
		if err != nil {
			return
		}
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) && opErr.Timeout() {
				continue
			}
			slog.Error("read udp error", err.Error(), slog.String("remoteAddr", remoteAddr.IP.String()))
		}

		if n < entity.MinRawPacketLength {
			slog.Warn("HepServerListener less then MinRawPacketLength",
				slog.Int("setting_length", entity.MinRawPacketLength),
				slog.Int("received_length", n),
				slog.String("remote_addr", remoteAddr.IP.String()),
			)
			continue
		}

		raw = make([]byte, n)

		copy(raw, data[:n])

		go ParseSaveOld(raw, remoteAddr.IP)
	}
}

func ParseSaveOld(b []byte, ip net.IP) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err != nil {
			slog.Error("parse save err", slog.String("msg", err.(error).Error()))
		}
	}()

	hepMsg, err := hep.NewHepMsg(b)

	if err != nil {
		slog.Error("NewHepMsg error", slog.String("err", err.Error()))
		return
	}

	if len(hepMsg.Body) <= 0 {
		return
	}

	if len(hepMsg.Body) < env.Conf.MinPacketLength {
		return
	}

	sip := parser.Parser{SIP: entity.SIP{}}
	sip.NodeIP = ip.String()
	sip.NodeID = strconv.Itoa(int(hepMsg.CaptureAgentID))

	bodyS := string(hepMsg.Body)

	sip.SIP.Raw = &bodyS
	sip.ParseCseq()

	sip.TimestampMicro = hepMsg.TimestampMicro

	if sip.CSeqMethod == "" {
		return
	}

	if strings.Contains(env.Conf.DiscardMethods, sip.CSeqMethod) {
		return
	}

	sip.ParseCallID()

	if sip.CallID == "" {
		return
	}

	sip.ParseFirstLine()

	if sip.Title == "" {
		return
	}

	if sip.RequestURL != "" {
		sip.ParseRequestURL()
	}

	sip.ParseFrom()
	sip.ParseTo()
	sip.ParseUserAgent()

	if len(sip.UserAgent) > entity.MaxUserAgentLength {
		sip.UserAgent = sip.UserAgent[:entity.MaxUserAgentLength]
	}

	sip.CreateTime = time.Unix(int64(hepMsg.Timestamp), 0)
	//s.TimestampMicro=s.CreateTime.Add(time.Microsecond * time.Duration(s.TimestampMicro)).UnixMicro()

	sip.Protocol = int(hepMsg.IPProtocolID)

	sip.SrcAddr = fmt.Sprintf("%s_%d", hepMsg.IP4SourceAddress, hepMsg.SourcePort)
	sip.SrcPort = int(hepMsg.SourcePort)
	sip.SrcHost = hepMsg.IP4SourceAddress
	sip.SrcCountryName, sip.SrcCityName, _ = GetIPArea(hepMsg.IP4SourceAddress)

	sip.DstAddr = fmt.Sprintf("%s_%d", hepMsg.IP4DestinationAddress, hepMsg.DestinationPort)
	sip.DstHost = hepMsg.IP4DestinationAddress
	sip.DstPort = int(hepMsg.DestinationPort)
	sip.DstCountryName, sip.DstCityName, _ = GetIPArea(hepMsg.IP4DestinationAddress)

	sip.UUID = fmt.Sprintf("%s%d%s", sip.NodeIP, hepMsg.CaptureAgentID, sip.CallID)

	model.SaveToDBQueue <- &sip.SIP
	return
}
