package services

import (
	"sip-monitor/src/config"
	"sip-monitor/src/model"

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
