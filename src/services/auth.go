package services

import (
	"sip-monitor/src/entity"

	"github.com/gin-gonic/gin"
)

func AuthLogin(c *gin.Context) {
	var request entity.AuthLogin
	err := c.ShouldBind(&request)
	if err != nil {
		return
	}

}
