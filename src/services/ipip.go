package services

import (
	"log/slog"

	"github.com/ipipdotnet/ipdb-go"
)

var IPDB *ipdb.City

func IPDBInit() {
	var err error
	IPDB, err = ipdb.NewCity("./ipv4.ipdb")
	if err != nil {
		slog.Error("IPAddressDatabaseInit Error", err.Error())
	} else {
		slog.Info("IPAddressDatabaseInit Success")
	}
}

func GetIPArea(ip string) (string, string) {
	if IPDB == nil {
		return "", ""
	}
	info, err := IPDB.FindInfo(ip, "CN")
	if err != nil {
		return "", ""
	}
	switch info.CountryName {
	case "局域网", "本机地址":
		return "局域网", "局域网"
	default:
		return info.CountryName, info.CityName
	}
}
