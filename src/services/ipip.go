package services

import (
	"log/slog"
	"sip-monitor/resources"

	"github.com/ipipdotnet/ipdb-go"
	"github.com/xiaoqidun/qqwry"
)

var IPDB *ipdb.City

func IPDBInit() {
	var err error
	IPDB, err = ipdb.NewCityFromBytes(resources.IPIPDat)
	if err != nil {
		slog.Error("IPAddressDatabaseInit ipip Error", slog.Any("err", err))
	} else {
		slog.Info("IPAddressDatabaseInit ipip Success")
	}
	qqwry.LoadData(resources.QQWryDat)
}

func GetIPArea(ip string) (string, string, string) {
	if IPDB == nil {
		return "", "", ""
	}
	info, err := IPDB.FindInfo(ip, "CN")
	if err != nil {
		return "", "", ""
	}
	switch info.CountryName {
	case "局域网", "本机地址":
		return "局域网", "局域网", ""
	default:
		return info.CountryName, info.CityName, info.IspDomain
	}
}

func GetIPAreaByCZ(ip string) (string, string, string) {
	address, err := qqwry.QueryIP(ip)
	if err != nil {
		return "", "", ""
	}

	return address.Country, address.Province + address.City, address.ISP
}
