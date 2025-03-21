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
	//查询出所有gateway, 并构建map，key为addr，value为name
	gatewayMap := make(map[string]string)
	gateways, _ := h.repository.GatewayList()
	for _, gateway := range gateways {
		gatewayMap[gateway.Addr] = gateway.Name
	}

	//将gatewayMap中的name替换为gatewayMap[stat.IP]
	for _, stat := range callStat {
		stat.Gateway = gatewayMap[stat.IP]
	}

	util.SendSuccessWithData(c, callStat)
}
