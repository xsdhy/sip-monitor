package hep

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/callbuffer"
	"sip-monitor/src/pkg/env"
	"sip-monitor/src/pkg/hep"
	"sip-monitor/src/pkg/parser"
	"sip-monitor/src/services/ip"
)

type HepServer struct {
	logger *logrus.Logger
	dal    model.DB
	conn   *net.UDPConn
	ip     *ip.IPServer

	saveQueue chan *entity.SIP

	callBufferMap map[string]*callbuffer.CallBuffer
}

func NewHepServer(logger *logrus.Logger, dal model.DB, ip *ip.IPServer) (*HepServer, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: env.Conf.UDPListenPort})
	if err != nil {
		logger.WithError(err).Error("HepServerListener Udp Service listen report udp fail")
		return nil, err
	}

	h := &HepServer{
		logger:        logger,
		dal:           dal,
		conn:          conn,
		ip:            ip,
		saveQueue:     make(chan *entity.SIP, 20000),
		callBufferMap: make(map[string]*callbuffer.CallBuffer),
	}
	return h, nil
}

func (h *HepServer) Listener() {
	defer h.conn.Close()
	h.logger.Info("HepServerListener")

	var data = make([]byte, env.Conf.MaxPacketLength)
	var raw []byte
	for {
		err := h.conn.SetDeadline(time.Now().Add(time.Duration(env.Conf.MaxReadTimeoutSeconds) * time.Second))
		if err != nil {
			return
		}
		n, remoteAddr, err := h.conn.ReadFromUDP(data)
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) && opErr.Timeout() {
				continue
			}
			h.logger.WithField("remoteAddr", remoteAddr.IP.String()).
				Error("read udp error", err.Error())
		}

		if n < entity.MinRawPacketLength {
			h.logger.
				WithField("setting_length", entity.MinRawPacketLength).
				WithField("received_length", n).
				WithField("remoteAddr", remoteAddr.IP.String()).
				Warn("HepServerListener less then MinRawPacketLength")
			continue
		}

		raw = make([]byte, n)

		copy(raw, data[:n])

		go func() {
			defer func() {
				// 发生宕机时，获取panic传递的上下文并打印
				err := recover()
				if err != nil {
					h.logger.WithError(err.(error)).Error("parse save err")
				}
			}()

			hepMsg, err := hep.NewHepMsg(raw)
			if err != nil {
				h.logger.WithError(err).Error("NewHepMsg error")
				return
			}

			item := h.parseSIP(hepMsg, remoteAddr.IP)
			if item == nil {
				return
			}
			h.saveQueue <- item
		}()
	}
}

func (h *HepServer) parseSIP(hepMsg *hep.HepMsg, ip net.IP) *entity.SIP {
	if len(hepMsg.Body) <= 0 || len(hepMsg.Body) < env.Conf.MinPacketLength {
		return nil
	}

	sip := parser.Parser{SIP: entity.SIP{}}

	sip.NodeIP = ip.String()
	sip.NodeID = strconv.Itoa(int(hepMsg.CaptureAgentID))
	sip.TimestampMicro = uint64(hepMsg.TimestampMicro)
	sip.CreateTime = time.Unix(int64(hepMsg.Timestamp), 0)
	sip.Protocol = int(hepMsg.IPProtocolID)

	bodyS := string(hepMsg.Body)

	sip.SIP.Raw = &bodyS
	sip.ParseCseq()

	if sip.CSeqMethod == "" {
		return nil
	}

	if strings.Contains(env.Conf.DiscardMethods, sip.CSeqMethod) {
		return nil
	}

	sip.ParseCallID()

	if sip.CallID == "" {
		return nil
	}

	sip.ParseFirstLine()

	if sip.Title == "" {
		return nil
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

	sip.SrcAddr = fmt.Sprintf("%s:%d", hepMsg.IP4SourceAddress, hepMsg.SourcePort)
	sip.SrcPort = int(hepMsg.SourcePort)
	sip.SrcHost = hepMsg.IP4SourceAddress

	sip.DstAddr = fmt.Sprintf("%s:%d", hepMsg.IP4DestinationAddress, hepMsg.DestinationPort)
	sip.DstHost = hepMsg.IP4DestinationAddress
	sip.DstPort = int(hepMsg.DestinationPort)

	if h.ip != nil {
		sip.SrcCountryName, sip.SrcCityName, _ = h.ip.GetIPArea(hepMsg.IP4SourceAddress)
		sip.DstCountryName, sip.DstCityName, _ = h.ip.GetIPArea(hepMsg.IP4DestinationAddress)
	}

	sip.UUID = fmt.Sprintf("%s%d%s", sip.NodeIP, hepMsg.CaptureAgentID, sip.CallID)

	return &sip.SIP
}

func (h *HepServer) SaveRunner() {
	for {
		select {
		case item := <-h.saveQueue:
			if item.CSeqMethod != "REGISTER" {
				h.dal.SaveMsg(item)
				h.saveCall(item)
			} else {
				//saveRegister(item)
			}
		}
	}
}

func (h *HepServer) saveCall(item *entity.SIP) {
	buffer, ok := h.callBufferMap[item.UUID]
	if !ok {
		buffer = callbuffer.NewCallBuffer(h.logger, h.saveItemCall)
		h.callBufferMap[item.UUID] = buffer
	}
	buffer.Add(item)
}
func (h *HepServer) saveItemCall(item *entity.SIPRecordCall) {
	if item.UUID != "" {
		h.dal.SaveCall(item)
	}
	delete(h.callBufferMap, item.UUID)
}
