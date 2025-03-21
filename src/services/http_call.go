package services

import (
	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/util"

	"github.com/gin-gonic/gin"
)

func (h *HandleHttp) CallList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.repository.GetCallList(c, request)
	util.SendItems(c, nil, records, meta)
}

func (h *HandleHttp) CallDetails(c *gin.Context) {
	sipCallID := c.Query("sip_call_id")

	callItem, err := h.repository.GetCallBySIPCallID(c, sipCallID)
	if err != nil {
		util.SendResponse(c, nil, entity.CallDetailsVO{})
		return
	}

	var vo entity.CallDetailsVO

	vo.Records = make([]entity.Record, 0)
	vo.Relevants = make([]entity.Record, 0)

	vo.Records, _ = h.repository.GetRecordsBySIPCallIDs(c, []string{sipCallID})

	if callItem.SessionID == "" {
		util.SendResponse(c, nil, vo)
		return
	}

	sipCallIDs, err := h.repository.GetCallIDsBySessionID(c, callItem.SessionID)
	if err != nil {
		util.SendResponse(c, nil, vo)
		return
	}
	vo.Relevants, _ = h.repository.GetRecordsBySIPCallIDs(c, sipCallIDs)

	util.SendResponse(c, nil, vo)
}

func (h *HandleHttp) RecordRaw(c *gin.Context) {
	idStr := c.Param("id")
	id, err := util.ParseInt64(idStr)
	if err != nil {
		util.SendError(c, err)
		return
	}
	recordRaw, err := h.repository.GetRecordRawByID(c, id)
	if err != nil {
		util.SendError(c, err)
		return
	}
	util.SendSuccessWithData(c, recordRaw)
}
