package env

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v10"
)

type config struct {
	UDPListenPort  int `env:"UDPListenPort" envDefault:"9060"`
	HTTPListenPort int `env:"HTTPListenPort" envDefault:"9059"`

	MaxPacketLength       int `env:"MaxPacketLength" envDefault:"4096"`
	MaxReadTimeoutSeconds int `env:"MaxReadTimeoutSecond" envDefault:"5"`

	DiscardMethods  string `env:"DiscardMethods" envDefault:"OPTIONS"`
	MinPacketLength int    `env:"MinPacketLength" envDefault:"24"`

	DBType     string `env:"DBType" envDefault:"file"`
	DSNURL     string `env:"DSN_URL" envDefault:""`
	DBUser     string `env:"DBUser" envDefault:""`
	DBPassword string `env:"DBPassword" envDefault:""`
	DBAddr     string `env:"DBAddr" envDefault:""`
	DBName     string `env:"DBName" envDefault:"call_sbc"`
}

var Conf = config{}

func init() {
	err := env.Parse(&Conf)
	if err != nil {
		slog.Error("%+v\n", err)
	}
	slog.Debug(fmt.Sprintf("%#v\n", Conf))
}
