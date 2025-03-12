package services

import (
	"time"

	"sip-monitor/src/model"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func Cron() {
	c := cron.New()
	_, err := c.AddFunc("@every 1h", DailyDelete)
	if err != nil {
		logrus.WithError(err).Error("定时任务添加失败")
	}
	c.Start()
}

func DailyDelete() {
	now := time.Now()
	day := now.Add(-1 * 24 * time.Hour)

	allRecordNum, err := model.CleanSipALL(&day)
	if err != nil {
		logrus.WithError(err).Error("DailyDelete CleanSipRecord error")
	}
	logrus.WithField("total", allRecordNum).Info("DailyDelete SipRecord ")
}
