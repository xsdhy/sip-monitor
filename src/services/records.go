package services

import (
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/util"

	"github.com/gin-gonic/gin"
)

func SearchCallID(c *gin.Context) {
	sipCallID := c.Query("call_id")
	records, _ := model.Infra.GetDetailsBySipCallID(c, entity.SearchParams{CallID: sipCallID})
	util.SendItems(c, nil, records, nil)
}

func RecordCallList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := model.Infra.GetRecordCallList(c, request)
	util.SendItems(c, nil, records, meta)
}
func RecordRegisterList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := model.Infra.GetRecordRegisterList(c, request)
	util.SendItems(c, nil, records, meta)
}
