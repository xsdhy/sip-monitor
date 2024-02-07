package services

import (
	"sbc/src/entity"
	"sbc/src/model"
	"time"

	"github.com/robfig/cron/v3"
)

func Cron() {
	c := cron.New()
	c.AddFunc("@daily", func() {
		now := time.Now()
		day5 := now.Add(-5 * 24 * time.Hour)
		day30 := now.Add(-30 * 24 * time.Hour)

		//删除30天前的所有记录
		other := entity.CleanSipRecordDTO{
			EndTime: &day30,
		}
		_, _ = model.CleanSipRecord(other)

		//删除5天前的所有注册记录
		register := entity.CleanSipRecordDTO{
			EndTime: &day5,
			Method:  "REGISTER",
		}
		_, _ = model.CleanSipRecord(register)
	})
	c.Start()
}
