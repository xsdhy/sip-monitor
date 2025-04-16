package rtcp

import (
	"encoding/json"
	"math"
	"testing"
)

func TestSummaryRTCPReport(t *testing.T) {
	// 测试用例1: 出站通话 c47f9f1f-9482-123e-cba3-fa163ebe21a3
	t.Run("出站通话数据测试", func(t *testing.T) {
		// 创建LegRTCPReport
		legReport := &LegRTCPReport{
			SrcAddr:    "185.32.76.185",
			SrcPort:    35781,
			DstAddr:    "192.168.11.141",
			DstPort:    18049,
			RawPackets: []*RTCPPacket{},
		}

		// 添加第一个RTCP包
		packet1 := `{"sender_information":{"ntp_timestamp_sec":3953699806,"ntp_timestamp_usec":0,"rtp_timestamp":3613598720,"packets":35,"octets":5600},"ssrc":1723189040,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":620378479,"fraction_lost":0,"packets_lost":0,"highest_seq_no":3543,"ia_jitter":63125527,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp1 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet1), rtcp1); err != nil {
			t.Fatalf("解析RTCP包1失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 添加第二个RTCP包
		packet2 := `{"sender_information":{"ntp_timestamp_sec":3953699811,"ntp_timestamp_usec":0,"rtp_timestamp":3613638720,"packets":219,"octets":34960},"ssrc":1723189040,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":620378479,"fraction_lost":0,"packets_lost":0,"highest_seq_no":3790,"ia_jitter":4294962322,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp2 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet2), rtcp2); err != nil {
			t.Fatalf("解析RTCP包2失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 添加第三个RTCP包
		packet3 := `{"sender_information":{"ntp_timestamp_sec":3953699816,"ntp_timestamp_usec":0,"rtp_timestamp":3613678720,"packets":373,"octets":59360},"ssrc":1723189040,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":620378479,"fraction_lost":0,"packets_lost":0,"highest_seq_no":4040,"ia_jitter":4294967295,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp3 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet3), rtcp3); err != nil {
			t.Fatalf("解析RTCP包3失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 断言结果
		if legReport.PacketLost != 0 {
			t.Errorf("总丢包数应为0，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 4040 {
			t.Errorf("总包数应为4040，但得到了 %d", legReport.PacketCount)
		}
		if legReport.PacketLostRate != 0 {
			t.Errorf("丢包率应为0%%，但得到了 %.2f%%", legReport.PacketLostRate)
		}

		// 打印计算结果，便于分析
		t.Logf("出站通话结果 - 总丢包: %d, 总包数: %d, 丢包率: %.2f%%, MOS: %.2f, 平均抖动: %d, 最大抖动: %d, 平均延迟: %d, 最大延迟: %d",
			legReport.PacketLost, legReport.PacketCount, legReport.PacketLostRate,
			legReport.Mos, legReport.JitterAvg, legReport.JitterMax, legReport.DelayAvg, legReport.DelayMax)
	})

	// 测试用例2: 入站通话 c4722c89-9482-123e-0e90-00163e0fafd1
	t.Run("入站通话数据测试", func(t *testing.T) {
		// 创建LegRTCPReport，以192.168.11.141:10457-47.87.10.92:10803方向的包为例
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.11.141",
			SrcPort:    10457,
			DstAddr:    "47.87.10.92",
			DstPort:    10803,
			RawPackets: []*RTCPPacket{},
		}

		// 添加第一个RTCP包
		packet1 := `{"sender_information":{"ntp_timestamp_sec":3953699809,"ntp_timestamp_usec":978728557,"rtp_timestamp":3613622960,"packets":125,"octets":19920},"ssrc":3870454559,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3959830815,"fraction_lost":1,"packets_lost":1,"highest_seq_no":37430,"ia_jitter":2635131,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3870454559}`
		rtcp1 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet1), rtcp1); err != nil {
			t.Fatalf("解析RTCP包1失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 添加第二个RTCP包
		packet2 := `{"sender_information":{"ntp_timestamp_sec":3953699813,"ntp_timestamp_usec":1064627903,"rtp_timestamp":3613655280,"packets":227,"octets":36160},"ssrc":3870454559,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3959830815,"fraction_lost":0,"packets_lost":1,"highest_seq_no":37631,"ia_jitter":6,"lsr":2883668950,"dlsr":254279}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3870454559}`
		rtcp2 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet2), rtcp2); err != nil {
			t.Fatalf("解析RTCP包2失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 添加第三个RTCP包
		packet3 := `{"sender_information":{"ntp_timestamp_sec":3953699817,"ntp_timestamp_usec":1150531544,"rtp_timestamp":3613677760,"packets":354,"octets":56320},"ssrc":3870454559,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3959830815,"fraction_lost":0,"packets_lost":1,"highest_seq_no":37815,"ia_jitter":71822628,"lsr":2883932404,"dlsr":254279}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3870454559}`
		rtcp3 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet3), rtcp3); err != nil {
			t.Fatalf("解析RTCP包3失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 断言结果
		if legReport.PacketLost != 1 { // 包的累计丢包数都是1
			t.Errorf("总丢包数应为1，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 37815 {
			t.Errorf("总包数应为37815，但得到了 %d", legReport.PacketCount)
		}
		expectedLossRate := float64(1) / float64(37815) * 100
		if legReport.PacketLostRate < expectedLossRate*0.95 || legReport.PacketLostRate > expectedLossRate*1.05 {
			t.Errorf("丢包率应约为%.4f%%，但得到了 %.4f%%", expectedLossRate, legReport.PacketLostRate)
		}

		// 打印计算结果，便于分析
		t.Logf("入站通话结果 - 总丢包: %d, 总包数: %d, 丢包率: %.4f%%, MOS: %.2f, 平均抖动: %d, 最大抖动: %d, 平均延迟: %d, 最大延迟: %d",
			legReport.PacketLost, legReport.PacketCount, legReport.PacketLostRate,
			legReport.Mos, legReport.JitterAvg, legReport.JitterMax, legReport.DelayAvg, legReport.DelayMax)
	})

	// 测试用例3: 另一个方向的入站通话
	t.Run("入站通话反向数据测试", func(t *testing.T) {
		// 创建LegRTCPReport，以47.87.10.92:10803-192.168.11.141:10457方向的包为例
		legReport := &LegRTCPReport{
			SrcAddr:    "47.87.10.92",
			SrcPort:    10803,
			DstAddr:    "192.168.11.141",
			DstPort:    10457,
			RawPackets: []*RTCPPacket{},
		}

		// 添加第一个RTCP包
		packet1 := `{"sender_information":{"ntp_timestamp_sec":3953699809,"ntp_timestamp_usec":1272345406,"rtp_timestamp":2987508238,"packets":176,"octets":28160},"ssrc":3959830815,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3870454559,"fraction_lost":2,"packets_lost":1,"highest_seq_no":21437,"ia_jitter":14,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3959830815}`
		rtcp1 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet1), rtcp1); err != nil {
			t.Fatalf("解析RTCP包1失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 添加第二个RTCP包
		packet2 := `{"sender_information":{"ntp_timestamp_sec":3953699813,"ntp_timestamp_usec":1358180328,"rtp_timestamp":2987540398,"packets":377,"octets":60320},"ssrc":3959830815,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3870454559,"fraction_lost":0,"packets_lost":1,"highest_seq_no":21539,"ia_jitter":7,"lsr":2883927924,"dlsr":2622}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3959830815}`
		rtcp2 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet2), rtcp2); err != nil {
			t.Fatalf("解析RTCP包2失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 添加第三个RTCP包
		packet3 := `{"sender_information":{"ntp_timestamp_sec":3953699817,"ntp_timestamp_usec":1444109738,"rtp_timestamp":92480,"packets":554,"octets":88640},"ssrc":3959830815,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3870454559,"fraction_lost":0,"packets_lost":1,"highest_seq_no":21665,"ia_jitter":0,"lsr":2884191379,"dlsr":2622}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3959830815}`
		rtcp3 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet3), rtcp3); err != nil {
			t.Fatalf("解析RTCP包3失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 断言结果
		if legReport.PacketLost != 1 { // 每个包的累计丢包数都是1
			t.Errorf("总丢包数应为1，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 21665 {
			t.Errorf("总包数应为21665，但得到了 %d", legReport.PacketCount)
		}

		// 检查抖动
		expectedJitterAvg := uint64(10) // (14+7+0)/2 四舍五入，因为计算方式不同
		if legReport.JitterAvg != expectedJitterAvg {
			t.Errorf("平均抖动应为%d，但得到了 %d", expectedJitterAvg, legReport.JitterAvg)
		}

		// 检查延迟
		expectedDelayAvg := uint64(2622) // 只有两个延迟值是2622，一个是0，所以平均是2622
		if legReport.DelayAvg != expectedDelayAvg {
			t.Errorf("平均延迟应为%d，但得到了 %d", expectedDelayAvg, legReport.DelayAvg)
		}

		// 打印计算结果，便于分析
		t.Logf("入站反向通话结果 - 总丢包: %d, 总包数: %d, 丢包率: %.4f%%, MOS: %.2f, 平均抖动: %d, 最大抖动: %d, 平均延迟: %d, 最大延迟: %d",
			legReport.PacketLost, legReport.PacketCount, legReport.PacketLostRate,
			legReport.Mos, legReport.JitterAvg, legReport.JitterMax, legReport.DelayAvg, legReport.DelayMax)
	})
}

// TestSummaryRTCPReportEdgeCases 测试边界条件和特殊情况
func TestSummaryRTCPReportEdgeCases(t *testing.T) {
	// 测试用例1: 空输入
	t.Run("空RTCP报告测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果都是默认值
		if legReport.PacketLost != 0 {
			t.Errorf("空输入时丢包数应为0，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 0 {
			t.Errorf("空输入时包数应为0，但得到了 %d", legReport.PacketCount)
		}
		if legReport.PacketLostRate != 0 {
			t.Errorf("空输入时丢包率应为0，但得到了 %.4f", legReport.PacketLostRate)
		}
		if legReport.Mos != 0 {
			t.Errorf("空输入时MOS应为0，但得到了 %.2f", legReport.Mos)
		}
	})

	// 测试用例2: 只有一个RTCP包
	t.Run("单个RTCP包测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// 添加一个SR包
		packet := `{"sender_information":{"ntp_timestamp_sec":123456,"ntp_timestamp_usec":789,"rtp_timestamp":101112,"packets":100,"octets":16000},"ssrc":1234567890,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":987654321,"fraction_lost":25,"packets_lost":10,"highest_seq_no":1000,"ia_jitter":30,"lsr":0,"dlsr":1500}],"report_blocks_xr":null,"sdes_ssrc":0}`
		rtcp := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet), rtcp); err != nil {
			t.Fatalf("解析RTCP包失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 10 {
			t.Errorf("单包时丢包数应为10，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 1000 {
			t.Errorf("单包时包数应为1000，但得到了 %d", legReport.PacketCount)
		}
		expectedLossRate := float64(10) / float64(1000) * 100
		if legReport.PacketLostRate != expectedLossRate {
			t.Errorf("单包时丢包率应为%.4f，但得到了 %.4f", expectedLossRate, legReport.PacketLostRate)
		}
		if legReport.JitterAvg != 30 {
			t.Errorf("单包时平均抖动应为30，但得到了 %d", legReport.JitterAvg)
		}
		if legReport.DelayAvg != 1500 {
			t.Errorf("单包时平均延迟应为1500，但得到了 %d", legReport.DelayAvg)
		}
	})

	// 测试用例3: 包含不支持的类型
	t.Run("不支持的RTCP类型测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// 创建正常的RR包
		packet1 := `{"ssrc":1234567890,"type":201,"report_count":1,"report_blocks":[{"source_ssrc":987654321,"fraction_lost":0,"packets_lost":5,"highest_seq_no":500,"ia_jitter":20,"lsr":0,"dlsr":1000}]}`
		rtcp1 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet1), rtcp1); err != nil {
			t.Fatalf("解析RTCP包1失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 创建不支持的BYE类型包 (203)
		packet2 := `{"ssrc":1234567890,"type":203}`
		rtcp2 := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet2), rtcp2); err != nil {
			t.Fatalf("解析RTCP包2失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 即使有不支持的包类型，也应该能处理支持的类型
		if legReport.PacketLost != 5 {
			t.Errorf("丢包数应为5，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 500 {
			t.Errorf("包数应为500，但得到了 %d", legReport.PacketCount)
		}
	})
}

// TestSummaryRTCPReportAbnormalData 测试异常数据处理
func TestSummaryRTCPReportAbnormalData(t *testing.T) {
	// 测试用例1: 异常大的抖动值
	t.Run("异常大抖动值测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// 正常抖动值
		packet1 := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":0,"packets_lost":0,"highest_seq_no":100,"ia_jitter":50,"lsr":0,"dlsr":0}]}`
		rtcp1 := &RTCPPacket{}
		json.Unmarshal([]byte(packet1), rtcp1)
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 异常大的抖动值
		packet2 := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":0,"packets_lost":0,"highest_seq_no":200,"ia_jitter":4294967290,"lsr":0,"dlsr":0}]}`
		rtcp2 := &RTCPPacket{}
		json.Unmarshal([]byte(packet2), rtcp2)
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 另一个正常值
		packet3 := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":0,"packets_lost":0,"highest_seq_no":300,"ia_jitter":70,"lsr":0,"dlsr":0}]}`
		rtcp3 := &RTCPPacket{}
		json.Unmarshal([]byte(packet3), rtcp3)
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证异常大的抖动值被忽略
		if legReport.JitterAvg != 60 { // (50+70)/2 = 60
			t.Errorf("过滤异常值后平均抖动应为60，但得到了 %d", legReport.JitterAvg)
		}
		if legReport.JitterMax != 70 {
			t.Errorf("过滤异常值后最大抖动应为70，但得到了 %d", legReport.JitterMax)
		}
	})

	// 测试用例2: 100%丢包率
	t.Run("极高丢包率测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// 100%丢包率
		packet := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":255,"packets_lost":1000,"highest_seq_no":1000,"ia_jitter":50,"lsr":0,"dlsr":0}]}`
		rtcp := &RTCPPacket{}
		json.Unmarshal([]byte(packet), rtcp)
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证100%丢包率能正确处理
		if legReport.PacketLostRate != 100.0 {
			t.Errorf("极端丢包率应为100.0%%，但得到了 %.4f%%", legReport.PacketLostRate)
		}
		if legReport.Mos < 1.0 || legReport.Mos > 1.1 { // 预期最低MOS值
			t.Errorf("极端丢包率时MOS应接近1.0，但得到了 %.2f", legReport.Mos)
		}
	})
}

// TestSummaryRTCPReportSpecialCases 测试特殊组合情况
func TestSummaryRTCPReportSpecialCases(t *testing.T) {
	// 测试用例1: 混合SR、RR和XR报告
	t.Run("混合报告类型测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// SR报告
		packet1 := `{"sender_information":{"ntp_timestamp_sec":123456,"ntp_timestamp_usec":789,"rtp_timestamp":101112,"packets":50,"octets":8000},"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":10,"packets_lost":5,"highest_seq_no":100,"ia_jitter":20,"lsr":0,"dlsr":500}]}`
		rtcp1 := &RTCPPacket{}
		json.Unmarshal([]byte(packet1), rtcp1)
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// RR报告
		packet2 := `{"ssrc":2222,"type":201,"report_count":1,"report_blocks":[{"source_ssrc":1111,"fraction_lost":5,"packets_lost":2,"highest_seq_no":150,"ia_jitter":30,"lsr":0,"dlsr":600}]}`
		rtcp2 := &RTCPPacket{}
		json.Unmarshal([]byte(packet2), rtcp2)
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// XR报告
		packet3 := `{"ssrc":3333,"type":207,"report_blocks_xr":{"type":7,"id":1,"fraction_lost":20,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":800,"end_system_delay":0}}`
		rtcp3 := &RTCPPacket{}
		json.Unmarshal([]byte(packet3), rtcp3)
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证混合报告能正确处理
		if legReport.PacketLost != 5 { // 取最大值5
			t.Errorf("混合报告丢包数应为5，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 150 {
			t.Errorf("混合报告包数应为150，但得到了 %d", legReport.PacketCount)
		}
		if legReport.JitterAvg != 25 { // (20 + 30) / 2 = 25
			t.Errorf("混合报告平均抖动应为25，但得到了 %d", legReport.JitterAvg)
		}
		if legReport.DelayAvg != 633 { // (500 + 600 + 800) / 3 = 633
			t.Errorf("混合报告平均延迟应为633，但得到了 %d", legReport.DelayAvg)
		}
	})

	// 测试用例2: 多个时间点的报告，值逐渐变化
	t.Run("时序报告变化测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.1.1",
			SrcPort:    10000,
			DstAddr:    "192.168.1.2",
			DstPort:    20000,
			RawPackets: []*RTCPPacket{},
		}

		// 第1个时间点：通话开始
		packet1 := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":0,"packets_lost":0,"highest_seq_no":100,"ia_jitter":5,"lsr":0,"dlsr":100}]}`
		rtcp1 := &RTCPPacket{}
		json.Unmarshal([]byte(packet1), rtcp1)
		rtcp1.TimestampMicro = 1000000 // 1秒
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 第2个时间点：网络变差
		packet2 := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":25,"packets_lost":10,"highest_seq_no":500,"ia_jitter":50,"lsr":0,"dlsr":500}]}`
		rtcp2 := &RTCPPacket{}
		json.Unmarshal([]byte(packet2), rtcp2)
		rtcp2.TimestampMicro = 5000000 // 5秒
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 第3个时间点：网络恢复
		packet3 := `{"ssrc":1111,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":2222,"fraction_lost":5,"packets_lost":15,"highest_seq_no":1000,"ia_jitter":20,"lsr":0,"dlsr":200}]}`
		rtcp3 := &RTCPPacket{}
		json.Unmarshal([]byte(packet3), rtcp3)
		rtcp3.TimestampMicro = 10000000 // 10秒
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证累积效果
		if legReport.PacketLost != 15 { // 最后一个报告中的累积丢包数
			t.Errorf("时序变化丢包数应为15，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 1000 {
			t.Errorf("时序变化包数应为1000，但得到了 %d", legReport.PacketCount)
		}
		expectedLossRate := float64(15) / float64(1000) * 100
		if math.Abs(legReport.PacketLostRate-expectedLossRate) > 0.001 {
			t.Errorf("时序变化丢包率应为%.4f%%，但得到了 %.4f%%", expectedLossRate, legReport.PacketLostRate)
		}
		if legReport.JitterAvg != 25 { // (5 + 50 + 20) / 3 = 25
			t.Errorf("时序变化平均抖动应为25，但得到了 %d", legReport.JitterAvg)
		}
	})
}

// TestMOSCalculation 专门测试MOS计算的正确性
func TestMOSCalculation(t *testing.T) {
	// 运行多个MOS计算测试，从丢包率0%到高丢包率
	testCases := []struct {
		name             string
		packetLost       uint64
		packetCnt        uint64
		jitter           uint64
		delay            uint64
		expectedMosRange []float64 // [min, max]
	}{
		{"完美通话质量", 0, 1000, 10, 100, []float64{4.3, 4.5}},
		{"良好通话质量", 5, 1000, 50, 200, []float64{3.5, 4.5}},
		{"一般通话质量", 50, 1000, 100, 300, []float64{3.8, 4.3}},
		{"较差通话质量", 150, 1000, 200, 500, []float64{2.0, 3.0}},
		{"极差通话质量", 300, 1000, 500, 1000, []float64{1.0, 1.5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			legReport := &LegRTCPReport{
				SrcAddr:    "192.168.1.1",
				SrcPort:    10000,
				DstAddr:    "192.168.1.2",
				DstPort:    20000,
				RawPackets: []*RTCPPacket{},
			}

			// 创建一个带有指定特性的RR包
			rtcp := &RTCPPacket{
				SSRC:        1111,
				PacketType:  RTCPPacketTypeRR,
				ReportCount: 1,
				ReportBlocks: []ReportBlock{
					{
						SourceSSRC:   2222,
						FractionLost: 0,
						PacketsLost:  tc.packetLost,
						HighestSeqNo: tc.packetCnt,
						IAJitter:     tc.jitter,
						DLSR:         tc.delay,
					},
				},
			}
			legReport.RawPackets = append(legReport.RawPackets, rtcp)

			// 执行汇总计算
			legReport.SummaryRTCPReport()

			// 验证MOS在预期范围内
			if legReport.Mos < tc.expectedMosRange[0] || legReport.Mos > tc.expectedMosRange[1] {
				t.Errorf("%s: MOS值应在%.1f-%.1f之间，但得到了 %.2f",
					tc.name, tc.expectedMosRange[0], tc.expectedMosRange[1], legReport.Mos)
			}

			// 打印MOS结果，方便分析
			t.Logf("%s: 丢包率=%.2f%%, 抖动=%d, 延迟=%d, MOS=%.2f",
				tc.name, legReport.PacketLostRate, legReport.JitterAvg, legReport.DelayAvg, legReport.Mos)
		})
	}
}

// TestSummaryRTCPReportRealWorldData 测试真实世界的RTCP数据
func TestSummaryRTCPReportRealWorldData(t *testing.T) {
	// 测试用例1: 源SSRC为0的SR报文
	t.Run("源SSRC为0的SR报文", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "185.32.76.193",
			SrcPort:    65507,
			DstAddr:    "192.168.11.141",
			DstPort:    13379,
			RawPackets: []*RTCPPacket{},
		}

		// SR报文，report_blocks中source_ssrc为0
		packet := `{"sender_information":{"ntp_timestamp_sec":3953720957,"ntp_timestamp_usec":0,"rtp_timestamp":3574752800,"packets":224,"octets":35680},"ssrc":319576335,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":0,"fraction_lost":0,"packets_lost":0,"highest_seq_no":0,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet), rtcp); err != nil {
			t.Fatalf("解析RTCP包失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 0 {
			t.Errorf("丢包数应为0，但得到了 %d", legReport.PacketLost)
		}
		// 由于highest_seq_no为0，包数也应为0
		if legReport.PacketCount != 0 {
			t.Errorf("包数应为0，但得到了 %d", legReport.PacketCount)
		}
		// 抖动和延迟也为0
		if legReport.JitterAvg != 0 || legReport.DelayAvg != 0 {
			t.Errorf("抖动和延迟应为0，但得到了抖动=%d, 延迟=%d", legReport.JitterAvg, legReport.DelayAvg)
		}
	})

	// 测试用例2: 带有SDES类型的报文
	t.Run("带有SDES类型的RTCP报文", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.11.141",
			SrcPort:    14691,
			DstAddr:    "47.87.10.92",
			DstPort:    16225,
			RawPackets: []*RTCPPacket{},
		}

		// SDES报文
		packet := `{"sender_information":{"ntp_timestamp_sec":3953720957,"ntp_timestamp_usec":1322321646,"rtp_timestamp":218960,"packets":1343,"octets":214880},"ssrc":3836208408,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":0,"fraction_lost":0,"packets_lost":1,"highest_seq_no":0,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3836208408}`
		rtcp := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet), rtcp); err != nil {
			t.Fatalf("解析RTCP包失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果 - 特别是SDES类型能被正确处理
		if legReport.PacketLost != 1 {
			t.Errorf("丢包数应为1，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 0 {
			t.Errorf("包数应为0，但得到了 %d", legReport.PacketCount)
		}
	})

	// 测试用例3: 包含有效报告块的SR报文
	t.Run("包含有效报告块的SR报文", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "185.32.76.189",
			SrcPort:    44469,
			DstAddr:    "192.168.11.141",
			DstPort:    10295,
			RawPackets: []*RTCPPacket{},
		}

		// 包含有效报告块的SR报文
		packet := `{"sender_information":{"ntp_timestamp_sec":3953720957,"ntp_timestamp_usec":0,"rtp_timestamp":544484480,"packets":113,"octets":18000},"ssrc":1892480545,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":3234686618,"fraction_lost":0,"packets_lost":0,"highest_seq_no":3579,"ia_jitter":2,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet), rtcp); err != nil {
			t.Fatalf("解析RTCP包失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 0 {
			t.Errorf("丢包数应为0，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 3579 {
			t.Errorf("包数应为3579，但得到了 %d", legReport.PacketCount)
		}
		if legReport.JitterAvg != 2 {
			t.Errorf("平均抖动应为2，但得到了 %d", legReport.JitterAvg)
		}
	})

	// 测试用例4: 包含丢包的RR报文
	t.Run("包含丢包的RR报文", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "45.181.15.37",
			SrcPort:    27693,
			DstAddr:    "192.168.11.141",
			DstPort:    13439,
			RawPackets: []*RTCPPacket{},
		}

		// 包含丢包的RR报文
		packet := `{"sender_information":{"ntp_timestamp_sec":3953720957,"ntp_timestamp_usec":1962799919,"rtp_timestamp":4095077017,"packets":5656,"octets":904960},"ssrc":4094762834,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":2297372210,"fraction_lost":0,"packets_lost":10,"highest_seq_no":68986,"ia_jitter":13,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":4094762834}`
		rtcp := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet), rtcp); err != nil {
			t.Fatalf("解析RTCP包失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 10 {
			t.Errorf("丢包数应为10，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 68986 {
			t.Errorf("包数应为68986，但得到了 %d", legReport.PacketCount)
		}
		expectedLossRate := float64(10) / float64(68986) * 100
		if math.Abs(legReport.PacketLostRate-expectedLossRate) > 0.001 {
			t.Errorf("丢包率应为%.6f%%，但得到了 %.6f%%", expectedLossRate, legReport.PacketLostRate)
		}
		if legReport.JitterAvg != 13 {
			t.Errorf("平均抖动应为13，但得到了 %d", legReport.JitterAvg)
		}
	})

	// 测试用例5: 包含fraction_lost=255的报文（极端丢包率）
	t.Run("极端丢包率报文", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "192.168.11.141",
			SrcPort:    18453,
			DstAddr:    "47.87.10.92",
			DstPort:    16639,
			RawPackets: []*RTCPPacket{},
		}

		// 包含fraction_lost=255的报文
		packet := `{"sender_information":{"ntp_timestamp_sec":3953720958,"ntp_timestamp_usec":1408233877,"rtp_timestamp":3320507201,"packets":184,"octets":29440},"ssrc":3877677615,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":0,"fraction_lost":255,"packets_lost":1,"highest_seq_no":0,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3877677615}`
		rtcp := &RTCPPacket{}
		if err := json.Unmarshal([]byte(packet), rtcp); err != nil {
			t.Fatalf("解析RTCP包失败: %v", err)
		}
		legReport.RawPackets = append(legReport.RawPackets, rtcp)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 1 {
			t.Errorf("丢包数应为1，但得到了 %d", legReport.PacketLost)
		}
		// 由于highest_seq_no为0，包数也为0，丢包率计算需要特殊处理
		// 这种情况下应该采用默认最差丢包率
		if legReport.PacketLostRate != 100.0 {
			t.Errorf("极端丢包率情况下丢包率应为100%%，但得到了 %.2f%%", legReport.PacketLostRate)
		}
		if legReport.Mos != 1.0 {
			t.Errorf("极端丢包率情况下MOS值应为1.0，但得到了 %.2f", legReport.Mos)
		}
	})

	// 测试用例6: 多个时序报告
	t.Run("真实多包RTCP时序报告", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "47.87.10.92",
			SrcPort:    13651,
			DstAddr:    "192.168.11.141",
			DstPort:    15663,
			RawPackets: []*RTCPPacket{},
		}

		// 第一个报告
		packet1 := `{"sender_information":{"ntp_timestamp_sec":3953720958,"ntp_timestamp_usec":413274638,"rtp_timestamp":1538240,"packets":7817,"octets":1250665},"ssrc":403473842,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":3835403330,"fraction_lost":0,"packets_lost":1,"highest_seq_no":18083,"ia_jitter":0,"lsr":4269625041,"dlsr":49808}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":403473842}`
		rtcp1 := &RTCPPacket{}
		json.Unmarshal([]byte(packet1), rtcp1)
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 另一个方向的对应报告
		packet2 := `{"sender_information":{"ntp_timestamp_sec":3953720958,"ntp_timestamp_usec":1408233877,"rtp_timestamp":3394159903,"packets":9577,"octets":1532320},"ssrc":3835403330,"type":202,"report_count":1,"report_blocks":[{"source_ssrc":403473842,"fraction_lost":0,"packets_lost":5,"highest_seq_no":67121,"ia_jitter":0,"lsr":4269676706,"dlsr":11796}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":3835403330}`
		rtcp2 := &RTCPPacket{}
		json.Unmarshal([]byte(packet2), rtcp2)
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 5 { // 取最大值5
			t.Errorf("丢包数应取最大值5，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 67121 {
			t.Errorf("包数应为67121，但得到了 %d", legReport.PacketCount)
		}
		expectedLossRate := float64(5) / float64(67121) * 100
		if math.Abs(legReport.PacketLostRate-expectedLossRate) > 0.001 {
			t.Errorf("丢包率应为%.6f%%，但得到了 %.6f%%", expectedLossRate, legReport.PacketLostRate)
		}

		// 延迟检查
		expectedDelayAvg := uint64((49808 + 11796) / 2)
		if legReport.DelayAvg != expectedDelayAvg {
			t.Errorf("平均延迟应为%d，但得到了 %d", expectedDelayAvg, legReport.DelayAvg)
		}
	})

	// 测试用例7: 混合多个真实世界RTCP包的组合测试
	t.Run("真实多包混合RTCP测试", func(t *testing.T) {
		legReport := &LegRTCPReport{
			SrcAddr:    "185.32.76.194",
			SrcPort:    39043,
			DstAddr:    "192.168.11.141",
			DstPort:    13673,
			RawPackets: []*RTCPPacket{},
		}

		// 添加SR包
		packet1 := `{"sender_information":{"ntp_timestamp_sec":3953720957,"ntp_timestamp_usec":0,"rtp_timestamp":4212693360,"packets":627,"octets":99920},"ssrc":2299398769,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":0,"fraction_lost":0,"packets_lost":0,"highest_seq_no":0,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp1 := &RTCPPacket{}
		json.Unmarshal([]byte(packet1), rtcp1)
		legReport.RawPackets = append(legReport.RawPackets, rtcp1)

		// 添加另一个SR包带有非零SSRC
		packet2 := `{"sender_information":{"ntp_timestamp_sec":3953720957,"ntp_timestamp_usec":0,"rtp_timestamp":2111573440,"packets":321,"octets":51280},"ssrc":2108949080,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":1291612225,"fraction_lost":0,"packets_lost":0,"highest_seq_no":53184,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp2 := &RTCPPacket{}
		json.Unmarshal([]byte(packet2), rtcp2)
		legReport.RawPackets = append(legReport.RawPackets, rtcp2)

		// 添加SR包带有非零丢包
		packet3 := `{"sender_information":{"ntp_timestamp_sec":3953720958,"ntp_timestamp_usec":0,"rtp_timestamp":3101516640,"packets":514,"octets":80720},"ssrc":873282488,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":688733585,"fraction_lost":0,"packets_lost":2,"highest_seq_no":55864,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp3 := &RTCPPacket{}
		json.Unmarshal([]byte(packet3), rtcp3)
		legReport.RawPackets = append(legReport.RawPackets, rtcp3)

		// 添加RR包带有抖动值
		packet4 := `{"sender_information":{"ntp_timestamp_sec":3953720958,"ntp_timestamp_usec":0,"rtp_timestamp":1359813520,"packets":479,"octets":75920},"ssrc":304320252,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":3637529270,"fraction_lost":0,"packets_lost":0,"highest_seq_no":21609,"ia_jitter":17,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`
		rtcp4 := &RTCPPacket{}
		json.Unmarshal([]byte(packet4), rtcp4)
		legReport.RawPackets = append(legReport.RawPackets, rtcp4)

		// 执行汇总计算
		legReport.SummaryRTCPReport()

		// 验证结果
		if legReport.PacketLost != 2 { // 取最大丢包值
			t.Errorf("丢包数应为2，但得到了 %d", legReport.PacketLost)
		}
		if legReport.PacketCount != 55864 { // 取最高序列号
			t.Errorf("包数应为55864，但得到了 %d", legReport.PacketCount)
		}
		if legReport.JitterAvg != 17 { // 只有一个报告中有抖动值17
			t.Errorf("平均抖动应为17，但得到了 %d", legReport.JitterAvg)
		}

		// 打印结果，便于分析
		t.Logf("混合RTCP测试结果 - 总丢包: %d, 总包数: %d, 丢包率: %.6f%%, MOS: %.2f, 平均抖动: %d, 平均延迟: %d",
			legReport.PacketLost, legReport.PacketCount, legReport.PacketLostRate,
			legReport.Mos, legReport.JitterAvg, legReport.DelayAvg)
	})
}
