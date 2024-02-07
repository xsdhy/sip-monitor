package services

import (
	_ "embed"
	"sip-monitor/resources"
	"testing"

	"github.com/ipipdotnet/ipdb-go"
	"github.com/xiaoqidun/qqwry"
)

func TestGetIPArea(t *testing.T) {
	qqwry.LoadData(resources.QQWryDat)
	IPDB, _ = ipdb.NewCityFromBytes(resources.IPIPDat)

	tests := []struct {
		ip   string
		city string
		isp  string
	}{
		{ip: "127.0.0.1", city: "本机地址", isp: ""},
		{ip: "10.10.10.2", city: "局域网", isp: "IP"},
		{ip: "114.64.247.219", city: "北京市", isp: "北京时代互通电信技术有限公司"},
		{ip: "110.187.122.110", city: "达州市", isp: "电信"},
		{ip: "194.213.3.102", city: "英国", isp: ""},
		{ip: "174.142.205.21", city: "加拿大", isp: "魁北克省蒙特利尔市iWeb科技公司"},
		{ip: "162.14.115.150", city: "成都市", isp: "腾讯云"},
		{ip: "124.114.150.42", city: "西安市", isp: "电信"},
		{ip: "223.86.169.226", city: "凉山州", isp: "移动"},
		{ip: "111.123.26.69", city: "毕节地区", isp: "电信"},
		{ip: "139.207.62.169", city: "成都市", isp: "电信"},
	}
	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			_, city, isp := GetIPArea(tt.ip)
			if city != tt.city {
				t.Errorf("GetIPArea() city = %v, want %v", city, tt.city)
			}
			if isp != tt.isp {
				t.Errorf("GetIPArea() isp = %v, want %v", isp, tt.isp)
			}
		})
	}
}
