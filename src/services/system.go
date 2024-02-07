package services

import (
	"fmt"
	"sbc/src/entity"
	"sbc/src/model"
	"sbc/src/pkg/util"

	"github.com/gin-gonic/gin"
)

func CleanSipRecord(c *gin.Context) {
	var request entity.CleanSipRecordDTO
	_ = c.ShouldBind(&request)
	if request.EndTime == nil {
		util.SendMessage(c, "请选择时间")
		return
	}
	deleteResult, err := model.CleanSipRecord(request)
	if err != nil {
		util.SendMessage(c, err.Error())
		return
	}
	util.SendSuccessByMessage(c, fmt.Sprintf("删除成功，共计%d条", deleteResult))
}

func DbStats(c *gin.Context) {
	res := make([]entity.MongoDBStatsVO, 0, 3)

	statsCallRecords := model.DbStats(c, "call_records")
	if statsCallRecords != nil {
		res = append(res, *statsCallRecords)
	}

	statsCallRecordsCall := model.DbStats(c, "call_records_call")
	if statsCallRecords != nil {
		res = append(res, *statsCallRecordsCall)
	}

	statsCallRecordsRegister := model.DbStats(c, "call_records_register")
	if statsCallRecords != nil {
		res = append(res, *statsCallRecordsRegister)
	}

	util.SendResponse(c, nil, res)
}
