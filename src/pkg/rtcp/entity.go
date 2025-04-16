package rtcp

import "encoding/json"

// RTCPPacketType 表示 RTCP 报文的类型（Packet Type）
type RTCPPacketType uint8

const (
	RTCPPacketTypeSR    RTCPPacketType = 200 // Sender Report (SR) 发送者报告
	RTCPPacketTypeRR    RTCPPacketType = 201 // Receiver Report (RR) 接收者报告
	RTCPPacketTypeSDES  RTCPPacketType = 202 // Source Description (SDES) 源描述
	RTCPPacketTypeBYE   RTCPPacketType = 203 // Goodbye (BYE) 再见
	RTCPPacketTypeAPP   RTCPPacketType = 204 // Application-defined (APP) 应用定义
	RTCPPacketTypeRTPFB RTCPPacketType = 205 // Transport layer feedback (RTPFB) 传输层反馈消息
	RTCPPacketTypePSFB  RTCPPacketType = 206 // Payload-specific feedback (PSFB) 特定负载反馈消息
	RTCPPacketTypeXR    RTCPPacketType = 207 // Extended Report (XR) 扩展报告
)

// XRBlockType 表示 RTCP XR（Extended Reports）中报告块的类型
type XRBlockType uint8

const (
	// 列举常见的 XR 报告块类型值（RFC 3611）：
	XRBlockTypeUnknown           XRBlockType = 0 // 未知或未使用
	XRBlockTypeLossRLE           XRBlockType = 1 // Loss RLE Report 丢包报告
	XRBlockTypeDuplicateRLE      XRBlockType = 2 // Duplicate RLE Report 重复报告
	XRBlockTypePacketReceiptTime XRBlockType = 3 // Packet Receipt Times Report 数据包接收时间报告
	XRBlockTypeReceiverRefTime   XRBlockType = 4 // Receiver Reference Time Report 接收方参考时间报告
	XRBlockTypeDLRR              XRBlockType = 5 // DLRR Report 延迟自上次接收者报告的时间
	XRBlockTypeStatisticsSummary XRBlockType = 6 // Statistics Summary Report 统计摘要报告
	XRBlockTypeVoIPMetrics       XRBlockType = 7 // VoIP Metrics Report 语音指标报告
)

// RTCPPacket 是一个不完整的，但是实际用到的RTCP包
type RTCPPacket struct {
	SSRC           uint32             `json:"ssrc"`                         // 报文自身的 SSRC
	PacketType     RTCPPacketType     `json:"type"`                         // RTCP 类型
	SenderInfo     *SenderInformation `json:"sender_information,omitempty"` // SR 才有
	ReportCount    uint8              `json:"report_count,omitempty"`       // SR/RR 才有(似乎都是1)
	ReportBlocks   []ReportBlock      `json:"report_blocks,omitempty"`      // SR/RR 才有
	ReportBlocksXR *XRReportBlock     `json:"report_blocks_xr,omitempty"`   // XR 才有
	SDESSSRC       uint32             `json:"sdes_ssrc,omitempty"`          // SDES 块中的 SSRC

	Raw            string `json:"-"`                         // 原始数据
	TimestampMicro int64  `json:"timestamp_micro,omitempty"` // 时间戳（微秒）
}

func (p *RTCPPacket) String() string {
	// 转JSON字符串
	json, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(json)
}

// SenderInformation 表示 SR 中的 Sender Info 块
type SenderInformation struct {
	NTPTimestampSec  uint64 `json:"ntp_timestamp_sec"`  // NTP 时间戳（秒）
	NTPTimestampUsec uint64 `json:"ntp_timestamp_usec"` // NTP 时间戳（微秒）
	RTPTimestamp     uint64 `json:"rtp_timestamp"`      // RTP 时间戳
	Packets          uint64 `json:"packets"`            // 发送包数
	Octets           uint64 `json:"octets"`             // 发送字节数
}

// ReportBlock 表示 SR/RR 中的单个 Report Block
type ReportBlock struct {
	SourceSSRC   uint32 `json:"source_ssrc"`    // 源 SSRC（表明这份报告是关于哪个媒体流）
	FractionLost uint8  `json:"fraction_lost"`  // 丢包率（1/256 单位）（最近一个周期内丢包的比例）
	PacketsLost  uint64 `json:"packets_lost"`   // 累计丢包数
	HighestSeqNo uint64 `json:"highest_seq_no"` // 最高序列号
	IAJitter     uint64 `json:"ia_jitter"`      // 抖动(单位是 RTP 时间戳)
	LSR          uint64 `json:"lsr"`            // 最后 SR 中的中间字段
	DLSR         uint64 `json:"dlsr"`           // 从收到最后 SR 到发送本 RR 的延迟
}

// XRReportBlock 表示单个 XR 扩展报告块
type XRReportBlock struct {
	Type            uint8  `json:"type"`             // XR 报告块类型
	ID              uint64 `json:"id"`               // 块标识符
	FractionLost    uint64 `json:"fraction_lost"`    // 丢包率
	FractionDiscard uint64 `json:"fraction_discard"` // 丢弃率
	BurstDensity    uint64 `json:"burst_density"`    // 突发密度
	GapDensity      uint64 `json:"gap_density"`      // 间隙密度
	BurstDuration   uint64 `json:"burst_duration"`   // 突发持续
	GapDuration     uint64 `json:"gap_duration"`     // 间隙持续
	RoundTripDelay  uint64 `json:"round_trip_delay"` // 往返时延
	EndSystemDelay  uint64 `json:"end_system_delay"` // 系统延迟
}
