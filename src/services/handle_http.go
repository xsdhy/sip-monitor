package services

import (
	"fmt"
	"sip-monitor/src/config"
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/hep"
	"sip-monitor/src/pkg/parser"
	"sip-monitor/src/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandleHttp struct {
	logger     *logrus.Logger
	cfg        *config.Config
	repository model.Repository
}

func NewHandleHttp(logger *logrus.Logger, cfg *config.Config, repository model.Repository) *HandleHttp {
	return &HandleHttp{
		logger:     logger,
		cfg:        cfg,
		repository: repository,
	}
}
func (h *HandleHttp) AuthLogin(c *gin.Context) {
	var request entity.AuthLogin
	err := c.ShouldBind(&request)
	if err != nil {
		return
	}

}

func (h *HandleHttp) RecordCallList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.repository.GetSIPCallRecordList(c, request)
	util.SendItems(c, nil, records, meta)
}

func (h *HandleHttp) SearchCallID(c *gin.Context) {
	sipCallID := c.Query("sip_call_id")
	var messages []entity.SIP
	records, _ := h.repository.GetRecordsBySIPCallID(c, sipCallID)
	for _, record := range records {
		hepMsg := hep.NewMockHepMsgBySIPMsg([]byte(record.RawMsg))
		parser := parser.NewParser(h.cfg, hepMsg)
		sip, err := parser.ParseSIPMsg()
		if err != nil {
			continue
		}
		sip.CreateAt = record.CreateTime
		sip.TimestampMicro = uint32(record.TimestampMicro)

		sip.SrcAddr = fmt.Sprintf("%s_%d", sip.FromDomain, sip.SrcPort)
		sip.DstAddr = strings.ReplaceAll(sip.ToDomain, ":", "_")

		messages = append(messages, *sip)
	}
	util.SendItems(c, nil, messages, nil)
}

func (h *HandleHttp) RecordRegisterList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.repository.GetSIPRegisterRecordList(c, request)
	util.SendItems(c, nil, records, meta)
}
