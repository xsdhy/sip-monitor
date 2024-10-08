package ip

import (
	"log/slog"

	"sip-monitor/resources"

	"github.com/ipipdotnet/ipdb-go"
	"github.com/xiaoqidun/qqwry"
)

type IPServer struct {
	db *ipdb.City
}

func NewIPServer() *IPServer {
	ipDb, err := ipdb.NewCityFromBytes(resources.IPIPDat)
	if err != nil {
		slog.Error("IPAddressDatabaseInit ipip Error", slog.Any("err", err))
	} else {
		slog.Info("IPAddressDatabaseInit ipip Success")
	}
	i := &IPServer{
		db: ipDb,
	}
	//qqwry.LoadData(resources.QQWryDat)
	return i
}

func (i *IPServer) GetIPArea(ip string) (string, string, string) {
	if i.db == nil {
		return "", "", ""
	}
	info, err := i.db.FindInfo(ip, "CN")
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

func (i *IPServer) GetIPAreaByCZ(ip string) (string, string, string) {
	address, err := qqwry.QueryIP(ip)
	if err != nil {
		return "", "", ""
	}

	return address.Country, address.Province + address.City, address.ISP
}
