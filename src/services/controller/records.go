package controller

import (
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/util"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	dal model.DB
}

func NewHttpServer(dal model.DB) *HttpServer {
	return &HttpServer{dal: dal}
}

func (h *HttpServer) SearchCallID(c *gin.Context) {
	sipCallID := c.Query("call_id")
	records, _ := h.dal.GetDetailsBySipCallID(c, entity.SearchParams{CallID: sipCallID})
	util.SendItems(c, nil, records, nil)
}

func (h *HttpServer) RecordCallList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.dal.GetRecordCallList(c, request)
	util.SendItems(c, nil, records, meta)
}
func (h *HttpServer) RecordRegisterList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.dal.GetRecordRegisterList(c, request)
	util.SendItems(c, nil, records, meta)
}
