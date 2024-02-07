package util

import (
	"net/http"
	"sbc/src/entity"
	"time"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := DecodeErr(err)
	nowTime := time.Now()
	// always return http.StatusOK
	c.JSON(http.StatusOK, entity.Response{Code: code, Msg: message, Data: data, Time: nowTime})
}

func SendItems(c *gin.Context, err error, data interface{}, meta *entity.Meta) {
	code, message := DecodeErr(err)
	nowTime := time.Now()
	// always return http.StatusOK
	c.JSON(http.StatusOK, entity.ResponseItems{Code: code, Msg: message, Data: data, Meta: meta, Time: nowTime})
}

func SendSuccess(c *gin.Context) {
	code, _ := DecodeErr(OK)
	nowTime := time.Now()
	// always return http.StatusOK
	c.JSON(http.StatusOK, entity.Response{Code: code, Msg: "ok", Data: nil, Time: nowTime})
}
func SendSuccessByMessage(c *gin.Context, msg string) {
	code, _ := DecodeErr(OK)
	nowTime := time.Now()
	// always return http.StatusOK
	c.JSON(http.StatusOK, entity.Response{Code: code, Msg: msg, Data: nil, Time: nowTime})
}
func SendMessage(c *gin.Context, message string) {
	code, _ := DecodeErr(ErrAll)
	nowTime := time.Now()
	// always return http.StatusOK
	c.JSON(http.StatusOK, entity.Response{Code: code, Msg: message, Data: nil, Time: nowTime})
}
