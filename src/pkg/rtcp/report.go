package rtcp

import "time"

type CallRTCPReports struct {
	CallID      string
	Legs        map[string]*LegRTCPReport
	LastUpdated time.Time // 最后更新时间
}

type LegRTCPReport struct {
	NodeIP  string // 节点IP
	SrcAddr string // 源地址
	SrcPort uint16 // 源端口
	DstAddr string // 目的地址
	DstPort uint16 // 目的端口

	Mos            float64 // 平均MOS
	PacketLost     uint64  // 总丢包数
	PacketCount    uint64  // 总包数
	PacketLostRate float64 // 丢包率
	JitterAvg      uint64  // 平均抖动
	JitterMax      uint64  // 抖动最大值
	DelayAvg       uint64  // 平均延迟
	DelayMax       uint64  // 延迟最大值

	RawPackets []*RTCPPacket // 通话过程中收到的RTCP包
}

// 根据（RawPackets []*RTCPPacket）汇总计算RTCP报告
// 汇总计算整个通话过程中的：MOS，丢包、抖动、延迟
func (s *LegRTCPReport) SummaryRTCPReport() {
	if len(s.RawPackets) == 0 {
		return
	}

	var totalPacketLost uint64
	var totalPackets uint64
	var totalJitter uint64
	var totalDelay uint64
	var jitterCount uint64
	var delayCount uint64
	var maxJitter uint64
	var maxDelay uint64
	var totalMos float64
	var mosCount int
	var maxPacketLost uint64
	var extremeLossDetected bool // 新增：检测极端丢包情况

	// 定义最大合理抖动和延迟阈值，超过这个值认为是异常
	const maxReasonableJitter uint64 = 10000000 // 10秒
	const maxReasonableDelay uint64 = 10000000  // 10秒

	// 遍历所有RTCP包
	for _, packet := range s.RawPackets {
		// 处理SR/RR报文
		if packet.PacketType == RTCPPacketTypeSR || packet.PacketType == RTCPPacketTypeRR || packet.PacketType == RTCPPacketTypeSDES {
			for _, block := range packet.ReportBlocks {
				// 检测极端丢包情况：fraction_lost=255或接近最大值
				if block.FractionLost >= 200 {
					extremeLossDetected = true
				}

				// 累计丢包数 - 取最大值
				if block.PacketsLost > maxPacketLost {
					maxPacketLost = block.PacketsLost
				}
				totalPacketLost += block.PacketsLost

				// 抖动统计 - 过滤异常大的值
				if block.IAJitter > 0 && block.IAJitter < maxReasonableJitter {
					totalJitter += uint64(block.IAJitter)
					jitterCount++
					if block.IAJitter > maxJitter {
						maxJitter = block.IAJitter
					}
				}

				// 延迟统计 (DLSR提供了延迟信息) - 过滤异常大的值
				if block.DLSR > 0 && block.DLSR < maxReasonableDelay {
					totalDelay += uint64(block.DLSR)
					delayCount++
					if block.DLSR > maxDelay {
						maxDelay = block.DLSR
					}
				}

				// 假设最高序列号可以作为总包数的估计 - 确保是递增的
				if block.HighestSeqNo > totalPackets {
					totalPackets = block.HighestSeqNo
				}
			}
		}

		// 处理XR扩展报告，特别关注VoIP指标
		if packet.PacketType == RTCPPacketTypeXR && packet.ReportBlocksXR != nil {
			xr := packet.ReportBlocksXR

			// XR块类型为VoIPMetrics(7)时包含MOS值
			if xr.Type == uint8(XRBlockTypeVoIPMetrics) {
				// 这里简化处理，实际MOS可能需要从特定字段提取
				// 假设我们通过其他指标计算或估算MOS
				// 例如使用R因子换算为MOS: MOS = 1 + 0.035*R + 7*10^-6*R*(R-60)*(100-R)
				r := 100.0 - (float64(xr.FractionLost) * 100.0 / 256.0)
				if r > 0 {
					mos := 1.0 + 0.035*r + 7e-6*r*(r-60)*(100-r)
					if mos > 0 && mos <= 5 { // MOS范围为1-5
						totalMos += mos
						mosCount++
					}
				}

				// 可能还有直接的延迟和抖动数据
				if xr.RoundTripDelay > 0 && xr.RoundTripDelay < maxReasonableDelay {
					totalDelay += uint64(xr.RoundTripDelay)
					delayCount++
					if xr.RoundTripDelay > maxDelay {
						maxDelay = xr.RoundTripDelay
					}
				}
			}
		}
	}

	// 计算平均值和更新结果
	s.PacketLost = maxPacketLost
	s.PacketCount = totalPackets

	// 计算丢包率 (避免除零错误)
	if totalPackets > 0 {
		s.PacketLostRate = float64(s.PacketLost) / float64(totalPackets) * 100
	} else if extremeLossDetected {
		// 特殊情况：检测到极端丢包但没有包数统计
		s.PacketLostRate = 100.0
	}

	// 平均抖动
	if jitterCount > 0 {
		s.JitterAvg = uint64(totalJitter / uint64(jitterCount))
	}
	s.JitterMax = maxJitter

	// 平均延迟
	if delayCount > 0 {
		s.DelayAvg = uint64(totalDelay / uint64(delayCount))
	}
	s.DelayMax = maxDelay

	// 平均MOS
	if mosCount > 0 {
		s.Mos = totalMos / float64(mosCount)
	} else {
		// 如果检测到极端丢包情况，直接设置为最低MOS值
		if extremeLossDetected {
			s.Mos = 1.0
		} else {
			// 使用改进的MOS计算公式，更符合E模型
			if s.PacketLostRate <= 100 {
				// R值计算 (简化的E-model)
				R := 93.2 // 基础值(无损情况)

				// 丢包影响 (根据丢包率影响R值)
				lossImpact := 0.0
				if s.PacketLostRate <= 2.0 {
					lossImpact = s.PacketLostRate * 0.5
				} else if s.PacketLostRate <= 10.0 {
					lossImpact = 1.0 + (s.PacketLostRate-2.0)*2.0
				} else {
					lossImpact = 17.0 + (s.PacketLostRate-10.0)*4.0
				}

				// 延迟影响 (根据延迟影响R值)
				delayImpact := 0.0
				if s.DelayAvg > 0 {
					delayMs := float64(s.DelayAvg) / 65.536 // 转换为毫秒
					delayImpact = 0.024 * delayMs
					if delayMs > 177.3 {
						delayImpact += 0.11 * (delayMs - 177.3)
					}
				}

				// 抖动影响 (根据抖动影响R值)
				jitterImpact := 0.0
				if s.JitterAvg > 0 {
					jitterMs := float64(s.JitterAvg)
					jitterImpact = jitterMs * 0.05
				}

				// 计算总的R值
				R = R - lossImpact - delayImpact - jitterImpact
				if R < 0 {
					R = 0
				}

				// 转换R值为MOS (ITU-T G.107)
				if R < 0 {
					s.Mos = 1.0
				} else if R > 100 {
					s.Mos = 4.5
				} else {
					s.Mos = 1 + 0.035*R + R*(R-60)*(100-R)*7e-6

					// 调整MOS值：让丢包率超过30%时，MOS值不会超过1.5
					if s.PacketLostRate > 30 {
						mosReduction := (s.PacketLostRate - 30) / 70 * 3.0 // 增加了降低比例
						if mosReduction > 3.0 {
							mosReduction = 3.0
						}
						s.Mos -= mosReduction
					}
				}
			} else {
				s.Mos = 1.0 // 最低MOS值
			}
		}

		// 确保MOS在1.0-5.0范围内
		if s.Mos < 1.0 {
			s.Mos = 1.0
		} else if s.Mos > 5.0 {
			s.Mos = 5.0
		}
	}
}
