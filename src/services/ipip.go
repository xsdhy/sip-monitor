package services

import (
	"fmt"
	"log/slog"
	"sip-monitor/resources"

	"github.com/ipipdotnet/ipdb-go"
	"github.com/pupuk/addr"
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
	address, isp, err := qqwry.QueryIP(ip)
	if err != nil {
		return "", "", ""
	}

	fmt.Println(address, isp)
	parse := addr.Smart(address)
	if parse.PostCode == "" {
		return address, "", isp
	}
	return "中国", parse.City, isp
}
