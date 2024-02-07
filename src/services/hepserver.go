package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"sbc/src/pkg/siprocket"
	"strconv"
	"strings"
	"time"

	"sbc/src/entity"
	"sbc/src/model"
	"sbc/src/pkg/env"
	"sbc/src/pkg/hep"
	"sbc/src/pkg/parser"
)

func HepServerListener() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: env.Conf.UDPListenPort})
	if err != nil {
		slog.Error("HepServerListener Udp Service listen report udp fail", err)
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
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else {
				slog.Error("read udp error", err, slog.String("remoteAddr", remoteAddr.IP.String()))
			}
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

func ParseSaveNew(b []byte, ip net.IP) *entity.Record {
	msg, err := hep.NewHepMsg(b)
	if err != nil {
		return nil
	}

	if msg.SipMsg == nil {
		msg.SipMsg = siprocket.Parse(msg.Body)
	}

	if msg.SipMsg == nil {
		slog.Warn("消息题为空", slog.String("ip", ip.String()), slog.Any("IPProtocolID", string(msg.IPProtocolID)))
		return nil
	}
	ua := msg.SipMsg.Ua.ToString()
	if len(ua) > entity.MaxUserAgentLength {
		ua = ua[:entity.MaxUserAgentLength]
	}
	item := entity.Record{
		ID: model.GetMd5(msg.SipMsg.CallId.ToString(), msg.SipMsg.Raw, ip.String()),

		NodeIP: ip.String(),

		SIPMethod:    string(msg.SipMsg.Req.Method),
		ResponseCode: BytesToInt(msg.SipMsg.Req.StatusCode),
		ResponseDesc: string(msg.SipMsg.Req.StatusDesc),

		CSeqMethod:  string(msg.SipMsg.Cseq.Method),
		CSeqNumber:  BytesToInt(msg.SipMsg.Cseq.Id),
		FromUser:    string(msg.SipMsg.From.User),
		FromHost:    string(msg.SipMsg.From.Host),
		ToUser:      string(msg.SipMsg.To.User),
		ToHost:      string(msg.SipMsg.To.Host),
		SIPCallID:   msg.SipMsg.CallId.ToString(),
		SIPProtocol: uint(msg.IPProtocolID),
		UserAgent:   ua,

		CreateTime:     time.Unix(int64(msg.Timestamp), 0),
		TimestampMicro: int64(msg.TimestampMicro),
		RawMsg:         msg.SipMsg.Raw,
	}

	item.SrcAddr = fmt.Sprintf("%s_%d", msg.SipMsg.From.Host, msg.SipMsg.From.Port)
	item.SrcHost = string(msg.SipMsg.From.Host)
	item.SrcPort = BytesToInt(msg.SipMsg.From.Port)
	item.SrcCountryName, item.SrcCityName = GetIPArea(item.SrcHost)

	item.DstAddr = fmt.Sprintf("%s_%d", msg.SipMsg.To.Host, msg.SipMsg.To.Port)
	item.DstHost = string(msg.SipMsg.From.Host)
	item.DstPort = BytesToInt(msg.SipMsg.From.Port)
	item.DstCountryName, item.DstCityName = GetIPArea(item.SrcHost)

	//model.Save(item,len(msg.SipMsg.Via))
	return &item
}

func ParseSaveOld(b []byte, ip net.IP) *entity.Record {
	s, errType, errMsg := Format(b)
	if errType != "" {
		slog.Warn("format msg error",
			slog.Int("raw_length", len(b)),
			slog.String("from", ip.String()),
			slog.String("err_type", errType),
			slog.String("err_msg", errMsg),
		)
		if errType != "method_discarded" {
			//slog.Error(fmt.Sprintf("format msg error: %v; raw length: %d, %s,  from: %v", errType, len(b), b, ip))
		}
		return nil
	}

	output := siprocket.Parse([]byte(*s.Raw))

	ua := s.UserAgent
	if len(ua) > entity.MaxUserAgentLength {
		ua = ua[:entity.MaxUserAgentLength]
	}
	item := entity.Record{
		ID:           model.GetMd5(s.CallID, *s.Raw, ip.String()),
		NodeIP:       ip.String(),
		FsCallID:     s.FSCallID,
		LegUid:       s.UID,
		SIPMethod:    s.Title,
		ResponseCode: s.ResponseCode,
		ResponseDesc: s.ResponseDesc,
		CSeqMethod:   s.CSeqMethod,
		CSeqNumber:   s.CSeqNumber,
		FromUser:     s.FromUsername,
		FromHost:     s.FromDomain,
		ToUser:       s.ToUsername,
		ToHost:       s.ToDomain,
		SIPCallID:    s.CallID,
		SIPProtocol:  uint(s.Protocol),
		UserAgent:    ua,

		SrcHost:        s.SrcHost,
		SrcPort:        s.SrcPort,
		SrcAddr:        s.SrcAddr,
		SrcCityName:    s.SrcCityName,
		SrcCountryName: s.SrcCountryName,

		DstHost:        s.DstHost,
		DstPort:        s.DstPort,
		DstAddr:        s.DstAddr,
		DstCityName:    s.DstCityName,
		DstCountryName: s.DstCountryName,

		CreateTime:     s.CreateAt,
		TimestampMicro: s.CreateAt.Add(time.Microsecond * time.Duration(s.TimestampMicro)).UnixMicro(),
		RawMsg:         *s.Raw,
	}

	model.Save(item, len(output.Via))
	return &item
}

func Format(p []byte) (s *entity.SIP, errorType string, errMsg string) {
	hepMsg, err := hep.NewHepMsg(p)

	if err != nil {
		slog.Error("NewHepMsg error %v", err)
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
	sip.SrcCountryName, sip.SrcCityName = GetIPArea(hepMsg.IP4SourceAddress)

	sip.DstAddr = fmt.Sprintf("%s_%d", hepMsg.IP4DestinationAddress, hepMsg.DestinationPort)
	sip.DstHost = hepMsg.IP4DestinationAddress
	sip.DstPort = int(hepMsg.DestinationPort)
	sip.DstCountryName, sip.DstCityName = GetIPArea(hepMsg.IP4DestinationAddress)

	sip.NodeID = strconv.Itoa(int(hepMsg.CaptureAgentID))

	return &sip.SIP, "", ""
}

func BytesToInt(bys []byte) int {
	var data int64
	_ = binary.Read(bytes.NewBuffer(bys), binary.BigEndian, &data)
	return int(data)
}
