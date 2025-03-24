package services

import (
	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/util"

	"github.com/gin-gonic/gin"
)

func (h *HandleHttp) CallStat(c *gin.Context) {
	var request entity.CallStatDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		util.SendError(c, err)
		return
	}

	callStat, err := h.repository.GetCallStat(c, request)
	if err != nil {
		util.SendError(c, err)
		return
	}
	util.SendSuccessWithData(c, callStat)
}
