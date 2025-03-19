package services

import (
	"sip-monitor/src/config"
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/siprocket"
	"sip-monitor/src/pkg/util"
	"strings"
	"time"

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

func (h *HandleHttp) RecordCallList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.repository.GetSIPCallRecordList(c, request)
	util.SendItems(c, nil, records, meta)
}

func (h *HandleHttp) CallDetails(c *gin.Context) {
	sipCallID := c.Query("sip_call_id")

	callItem, err := h.repository.GetSIPCallRecordBySIPCallID(c, sipCallID)
	if err != nil {
		util.SendResponse(c, nil, entity.CallDetailsVO{})
		return
	}

	var vo entity.CallDetailsVO

	vo.Records = make([]entity.SIP, 0)
	vo.Relevants = make([]entity.SIP, 0)

	records, _ := h.repository.GetRecordsBySIPCallID(c, sipCallID)
	for _, record := range records {
		sip := siprocket.ParseSIP([]byte(record.Raw))
		if sip == nil {
			continue
		}
		sip.CreateTime = record.CreateTime
		sip.TimestampMicro = record.TimestampMicro
		sip.SrcAddr = strings.ReplaceAll(record.SrcAddr, ":", "_")
		sip.DstAddr = strings.ReplaceAll(record.DstAddr, ":", "_")
		sip.Raw = &record.Raw
		vo.Records = append(vo.Records, *sip)
	}
	if callItem.SessionID == "" {
		util.SendResponse(c, nil, vo)
		return
	}

	sipCallIDs, err := h.repository.GetSIPCallIDsBySessionID(c, callItem.SessionID)
	if err != nil {
		util.SendResponse(c, nil, vo)
		return
	}

	relevantRecords, err := h.repository.GetRecordsBySIPCallIDs(c, sipCallIDs)
	if err != nil {
		util.SendResponse(c, nil, vo)
		return
	}

	for _, record := range relevantRecords {
		sip := siprocket.ParseSIP([]byte(record.Raw))
		if sip == nil {
			continue
		}
		sip.CreateTime = record.CreateTime
		sip.TimestampMicro = record.TimestampMicro
		sip.SrcAddr = strings.ReplaceAll(record.SrcAddr, ":", "_")
		sip.DstAddr = strings.ReplaceAll(record.DstAddr, ":", "_")
		sip.Raw = &record.Raw
		vo.Relevants = append(vo.Relevants, *sip)
	}
	util.SendResponse(c, nil, vo)
}

func (h *HandleHttp) RecordRegisterList(c *gin.Context) {
	var request entity.SearchParams
	_ = c.ShouldBind(&request)

	records, meta, _ := h.repository.GetSIPRegisterRecordList(c, request)
	util.SendItems(c, nil, records, meta)
}

func (h *HandleHttp) UserList(c *gin.Context) {
	var users []entity.User
	var err error

	// 从数据库获取所有用户
	users, err = h.repository.GetUsers(c)
	if err != nil {
		util.SendError(c, err)
		return
	}

	// 返回用户列表
	util.SendItems(c, nil, users, nil)
}

func (h *HandleHttp) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := util.ParseInt64(idStr)
	if err != nil {
		util.SendError(c, err)
		return
	}

	user, err := h.repository.GetUserByID(c, id)
	if err != nil {
		util.SendError(c, err)
		return
	}

	util.SendSuccessWithData(c, user)
}

func (h *HandleHttp) CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		util.SendError(c, err)
		return
	}

	// 设置创建和更新时间
	now := time.Now()
	user.CreateAt = now
	user.UpdateAt = now

	// 对密码进行加密
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.SendError(c, err)
		return
	}
	user.Password = hashedPassword

	if err := h.repository.CreateUser(c, &user); err != nil {
		util.SendError(c, err)
		return
	}

	util.SendSuccessWithData(c, user)
}

func (h *HandleHttp) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := util.ParseInt64(idStr)
	if err != nil {
		util.SendError(c, err)
		return
	}

	// 先获取原有用户信息
	existingUser, err := h.repository.GetUserByID(c, id)
	if err != nil {
		util.SendError(c, err)
		return
	}

	// 解析请求体
	var updateUser entity.User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		util.SendError(c, err)
		return
	}

	// 更新字段
	existingUser.Nickname = updateUser.Nickname
	existingUser.Username = updateUser.Username
	existingUser.UpdateAt = time.Now()

	// 如果提供了新密码，则更新密码
	if updateUser.Password != "" {
		hashedPassword, err := util.HashPassword(updateUser.Password)
		if err != nil {
			util.SendError(c, err)
			return
		}
		existingUser.Password = hashedPassword
	}

	// 保存更新
	if err := h.repository.UpdateUser(c, existingUser); err != nil {
		util.SendError(c, err)
		return
	}

	util.SendSuccessWithData(c, existingUser)
}

func (h *HandleHttp) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := util.ParseInt64(idStr)
	if err != nil {
		util.SendError(c, err)
		return
	}

	if err := h.repository.DeleteUser(c, id); err != nil {
		util.SendError(c, err)
		return
	}

	util.SendSuccess(c)
}
