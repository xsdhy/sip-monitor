package services

import (
	"fmt"
	"sip-monitor/resources"

	"github.com/ipipdotnet/ipdb-go"
	"github.com/pupuk/addr"
	"github.com/sirupsen/logrus"
	"github.com/xiaoqidun/qqwry"
)

var IPDB *ipdb.City

func IPDBInit() {
	var err error
	IPDB, err = ipdb.NewCityFromBytes(resources.IPIPDat)
	if err != nil {
		logrus.WithError(err).Error("IPAddressDatabaseInit ipip Error")
	} else {
		logrus.Info("IPAddressDatabaseInit ipip Success")
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
	location, err := qqwry.QueryIP(ip)
	if err != nil {
		return "", "", ""
	}

	fmt.Println(location.Country, location.Province, location.City, location.ISP)
	parse := addr.Smart(location.Country + location.Province + location.City)
	if parse.PostCode == "" {
		return location.Country + location.Province + location.City, "", location.ISP
	}
	return "中国", parse.City, location.ISP
}
