package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/sirupsen/logrus"
)

type Config struct {
	UDPListenPort  int `env:"UDPListenPort" envDefault:"9060"`
	HTTPListenPort int `env:"HTTPListenPort" envDefault:"9059"`

	MaxPacketLength       int `env:"MaxPacketLength" envDefault:"4096"`
	MaxReadTimeoutSeconds int `env:"MaxReadTimeoutSecond" envDefault:"5"`

	HeaderSessionIDName string `env:"HeaderSessionIDName" envDefault:"X-JCallId"`

	DiscardMethods  string `env:"DiscardMethods" envDefault:"OPTIONS"`
	MinPacketLength int    `env:"MinPacketLength" envDefault:"24"`

	DBType     string `env:"DBType" envDefault:"sqlite"`
	DSNURL     string `env:"DSN_URL" envDefault:""`
	DBUser     string `env:"DBUser" envDefault:""`
	DBPassword string `env:"DBPassword" envDefault:""`
	DBAddr     string `env:"DBAddr" envDefault:""`
	DBName     string `env:"DBName" envDefault:"monitor"`
	DBPath     string `env:"DBPath" envDefault:""`

	// JWT Authentication settings
	JWTSecret      string `env:"JWT_SECRET" envDefault:"sip-monitor-secret-key"`
	JWTExpiryHours int    `env:"JWT_EXPIRY_HOURS" envDefault:"1200"`
}

func ParseConfig() (Config, error) {
	var Conf Config
	err := env.Parse(&Conf)
	if err != nil {
		logrus.WithError(err).Error("env.Parse error")
	}
	logrus.Debugf("%#v\n", Conf)
	return Conf, nil
}
