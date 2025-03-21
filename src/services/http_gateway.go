package services

import (
	"errors"
	"strconv"
	"time"

	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *HandleHttp) GatewayCreate(c *gin.Context) {
	var req entity.Gateway
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, err)
		return
	}
	now := time.Now()
	gateway := &entity.Gateway{
		Name:     req.Name,
		Addr:     req.Addr,
		Remark:   req.Remark,
		CreateAt: &now,
		UpdateAt: &now,
	}
	h.repository.GatewayCreate(gateway)
	util.SendSuccess(c)
}

func (h *HandleHttp) GatewayUpdate(c *gin.Context) {
	var req entity.Gateway
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, err)
		return
	}
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.SendError(c, err)
		return
	}
	now := time.Now()
	var gateway entity.Gateway
	gateway.ID = idInt
	gateway.Name = req.Name
	gateway.Addr = req.Addr
	gateway.Remark = req.Remark
	gateway.UpdateAt = &now
	h.repository.GatewayUpdate(&gateway)
	util.SendSuccess(c)
}

func (h *HandleHttp) GatewayGetByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.SendError(c, err)
		return
	}
	gateway, err := h.repository.GatewayGetByID(idInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.SendError(c, err)
			return
		}
		util.SendError(c, err)
		return
	}
	util.SendSuccessWithData(c, gateway)
}

func (h *HandleHttp) GatewayList(c *gin.Context) {
	gateways, err := h.repository.GatewayList()
	if err != nil {
		util.SendError(c, err)
		return
	}

	util.SendSuccessWithData(c, gateways)
}

func (h *HandleHttp) GatewayDelete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.SendError(c, err)
		return
	}
	h.repository.GatewayDelete(idInt)
	util.SendSuccess(c)
}
