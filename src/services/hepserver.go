package services

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"sip-monitor/src/config"
	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/hep"
	"sip-monitor/src/pkg/siprocket"

	"github.com/sirupsen/logrus"
)

type HepServer struct {
	logger      *logrus.Logger
	conn        *net.UDPConn
	cfg         *config.Config
	saveService *SaveService
}

func NewHepServer(logger *logrus.Logger, cfg *config.Config, saveService *SaveService) (*HepServer, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: cfg.UDPListenPort})
	if err != nil {
		logger.WithError(err).Error("HepServerListener Udp Service listen report udp fail")
		return nil, err
	}
	return &HepServer{
		conn:        conn,
		logger:      logger,
		cfg:         cfg,
		saveService: saveService,
	}, nil
}

func (h *HepServer) Start() error {
	defer h.conn.Close()
	h.logger.Info("HepServerListener")

	if h.cfg.MaxPacketLength <= 0 {
		h.cfg.MaxPacketLength = 4096
	}
	var data = make([]byte, h.cfg.MaxPacketLength)
	var raw []byte
	for {
		err := h.conn.SetDeadline(time.Now().Add(time.Duration(h.cfg.MaxReadTimeoutSeconds) * time.Second))
		if err != nil {
			return err
		}
		n, remoteAddr, err := h.conn.ReadFromUDP(data)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else {
				h.logger.WithFields(logrus.Fields{
					"remote_addr": remoteAddr.IP.String(),
				}).WithError(err).Error("read udp error")
			}
		}

		if n < entity.MinRawPacketLength {
			h.logger.WithFields(logrus.Fields{
				"setting_length":  entity.MinRawPacketLength,
				"received_length": n,
				"remote_addr":     remoteAddr.IP.String(),
			}).Warn("HepServerListener less then MinRawPacketLength")
			continue
		}

		raw = make([]byte, n)

		copy(raw, data[:n])

		go h.ParseSIPMsg(raw, remoteAddr.IP.String())
	}
}

func (h *HepServer) ParseSIPMsg(b []byte, ip string) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err != nil {
			h.logger.WithError(err.(error)).Error("parse save err")
		}
	}()

	hepMsg, err := hep.NewHepMsg(b)
	if err != nil {
		return
	}
	if len(hepMsg.Body) <= 0 {
		return
	}

	if len(hepMsg.Body) < h.cfg.MinPacketLength {
		return
	}

	sip := siprocket.ParseSIP(hepMsg.Body)
	if sip == nil {
		return
	}

	sip.CreateTime = time.Unix(int64(hepMsg.Timestamp), 0)
	sip.TimestampMicro = sip.CreateTime.Add(time.Microsecond * time.Duration(hepMsg.TimestampMicro)).UnixMicro()

	sip.Protocol = int(hepMsg.IPProtocolID)

	sip.SrcAddr = fmt.Sprintf("%s:%d", hepMsg.IP4SourceAddress, hepMsg.SourcePort)
	sip.DstAddr = fmt.Sprintf("%s:%d", hepMsg.IP4DestinationAddress, hepMsg.DestinationPort)

	sip.NodeID = strconv.Itoa(int(hepMsg.CaptureAgentID))
	sip.NodeIP = ip

	h.saveService.SaveToDBQueue <- *sip
}
