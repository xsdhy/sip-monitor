package services

import (
	"sbc/src/entity"
	"sbc/src/model"
	"sbc/src/pkg/util"

	"github.com/gin-gonic/gin"
)

func SearchCallID(c *gin.Context) {
	sipCallID := c.Query("sip_call_id")
	records, _ := model.GetDetailsBySipCallID(c, entity.SearchParams{SipCallID: sipCallID})
	util.SendItems(c, nil, records, nil)
}

func SearchAll(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := model.GetRecordList(c, request)
	util.SendItems(c, nil, records, meta)
}

func RecordCallList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := model.GetRecordCallList(c, request)
	util.SendItems(c, nil, records, meta)
}
func RecordRegisterList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := model.GetRecordRegisterList(c, request)
	util.SendItems(c, nil, records, meta)
}
