package services

import (
	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandleHttp) UserList(c *gin.Context) {
	var users []entity.User
	var err error

	users, err = h.repository.GetUsers(c)
	if err != nil {
		util.SendError(c, err)
		return
	}

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

	now := time.Now()
	user.CreateAt = now
	user.UpdateAt = now

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
