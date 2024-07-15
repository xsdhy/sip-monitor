package services

import (
	"log/slog"
	"time"

	"sip-monitor/src/model"

	"github.com/robfig/cron/v3"
)

func Cron() {
	c := cron.New()
	_, err := c.AddFunc("@every 1h", DailyDelete)
	if err != nil {
		slog.Error("定时任务添加失败:", err)
	}
	c.Start()
}

func DailyDelete() {
	now := time.Now()
	day := now.Add(-1 * 24 * time.Hour)

	allRecordNum, err := model.CleanSipALL(&day)
	if err != nil {
		slog.Error("DailyDelete CleanSipRecord error:", err)
	}
	slog.Info("DailyDelete SipRecord ", slog.Int64("total", allRecordNum))
}
