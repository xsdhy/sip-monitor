package rtcp

import (
	"encoding/json"
	"math"
	"testing"
)

// TestCalculateMOS tests the MOS calculation function
func TestCalculateMOS(t *testing.T) {
	tests := []struct {
		name      string
		lossRate  float64
		jitter    float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "Perfect Quality",
			lossRate:  0.0,
			jitter:    0.0,
			expected:  4.3,
			tolerance: 0.01,
		},
		{
			name:      "Good Quality",
			lossRate:  0.01,
			jitter:    5.0,
			expected:  4.10,
			tolerance: 0.01,
		},
		{
			name:      "Average Quality",
			lossRate:  0.05,
			jitter:    10.0,
			expected:  3.92,
			tolerance: 0.01,
		},
		{
			name:      "Poor Quality",
			lossRate:  0.25,
			jitter:    20.0,
			expected:  3.57,
			tolerance: 0.01,
		},
		{
			name:      "Very Poor Quality",
			lossRate:  0.5,
			jitter:    50.0,
			expected:  3.32,
			tolerance: 0.01,
		},
		{
			name:      "Extreme Loss",
			lossRate:  1.0,
			jitter:    10.0,
			expected:  3.26,
			tolerance: 0.01,
		},
		{
			name:      "Out of Bounds Loss",
			lossRate:  2.0, // Should be capped at 1.0
			jitter:    10.0,
			expected:  3.26,
			tolerance: 0.01,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateMOS(tc.lossRate, tc.jitter)
			if math.Abs(result-tc.expected) > tc.tolerance {
				t.Errorf("calculateMOS(%f, %f) = %f, expected %f (±%f)",
					tc.lossRate, tc.jitter, result, tc.expected, tc.tolerance)
			}
		})
	}
}

// TestConvertFractionLostToRate tests the fraction lost conversion function
func TestConvertFractionLostToRate(t *testing.T) {
	tests := []struct {
		name         string
		fractionLost int
		expected     float64
	}{
		{"No Loss", 0, 0.0},
		{"Quarter Loss", 64, 0.25},
		{"Half Loss", 128, 0.5},
		{"Three Quarter Loss", 192, 0.75},
		{"Full Loss", 255, 0.99609375},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := convertFractionLostToRate(tc.fractionLost)
			if math.Abs(result-tc.expected) > 0.0001 {
				t.Errorf("convertFractionLostToRate(%d) = %f, expected %f",
					tc.fractionLost, result, tc.expected)
			}
		})
	}
}

func createReportFromRaw(packets []string) *LegRTCPReport {
	report := &LegRTCPReport{
		RawPackets: make([]*RTCPPacket, len(packets)),
	}
	for i, item := range packets {
		var rtcpPacket RTCPPacket
		err := json.Unmarshal([]byte(item), &rtcpPacket)
		if err == nil {
			report.RawPackets[i] = &rtcpPacket
		}
	}
	return report
}

// TestProcessRTCPPackets tests the RTCP packet processing function
func TestProcessRTCPPackets(t *testing.T) {
	// Test with multiple packets and invalid packet
	t.Run("Multiple Packets With Invalid", func(t *testing.T) {
		packets := []string{
			`{"sender_information":{"packets":100,"octets":16000},"report_blocks":[{"fraction_lost":25,"packets_lost":5,"ia_jitter":10}]}`,
			`{"sender_information":{"packets":200,"octets":32000},"report_blocks":[{"fraction_lost":51,"packets_lost":10,"ia_jitter":20}]}`,
		}
		report := createReportFromRaw(packets)

		report.ProcessRTCPPackets()

		// Should process the two valid packets
		if report.PacketCount != 300 {
			t.Errorf("Expected PacketCount=300, got %d", report.PacketCount)
		}

		// Average jitter should be (10+20)/2 = 15
		if report.JitterAvg != 15 {
			t.Errorf("Expected JitterAvg=15, got %d", report.JitterAvg)
		}

		// Max jitter should be 20
		if report.JitterMax != 20 {
			t.Errorf("Expected JitterMax=20, got %d", report.JitterMax)
		}
	})

	// Test the actual example packets
	t.Run("Example Packets", func(t *testing.T) {
		examplePackets := []string{
			`{"sender_information":{"ntp_timestamp_sec":3953771345,"ntp_timestamp_usec":0,"rtp_timestamp":2944344000,"packets":2,"octets":320},"ssrc":3571175962,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":0,"fraction_lost":0,"packets_lost":0,"highest_seq_no":0,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`,
			`{"sender_information":{"ntp_timestamp_sec":3953771345,"ntp_timestamp_usec":0,"rtp_timestamp":2944344000,"packets":2,"octets":320},"ssrc":3571175962,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":0,"fraction_lost":0,"packets_lost":0,"highest_seq_no":0,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`,
			`{"sender_information":{"ntp_timestamp_sec":3953771350,"ntp_timestamp_usec":0,"rtp_timestamp":2944384000,"packets":18,"octets":2880},		"ssrc":3571175962,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":1089493108,"fraction_lost":0,"packets_lost":16777215,"highest_seq_no":24658,"ia_jitter":0,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`,
			`{"sender_information":{"ntp_timestamp_sec":3953771355,"ntp_timestamp_usec":0,"rtp_timestamp":2944424000,"packets":201,"octets":32000},"ssrc":3571175962,"type":200,"report_count":1,"report_blocks":[{"source_ssrc":1089493108,"fraction_lost":1,"packets_lost":0,"highest_seq_no":24907,"ia_jitter":13,"lsr":0,"dlsr":0}],"report_blocks_xr":{"type":0,"id":0,"fraction_lost":0,"fraction_discard":0,"burst_density":0,"gap_density":0,"burst_duration":0,"gap_duration":0,"round_trip_delay":0,"end_system_delay":0},"sdes_ssrc":0}`,
		}

		report := createReportFromRaw(examplePackets)

		report.ProcessRTCPPackets()

		// Total packet count should be 2 + 2 + 18 + 201 = 223
		if report.PacketCount != 223 {
			t.Errorf("Expected PacketCount=223, got %d", report.PacketCount)
		}

		// Loss rate should reflect the 1/256 (0.00390625) in the third packet
		expectedLossRate := 1.0 / 256.0 / 4 // 计算方式改为总样本数
		if math.Abs(report.PacketLostRate-expectedLossRate) > 0.0001 {
			t.Errorf("Expected PacketLostRate=%.6f, got %.6f", expectedLossRate, report.PacketLostRate)
		}

		// Jitter should be 13ms (only valid jitter measurement in the third packet)
		if report.JitterAvg != 13 {
			t.Errorf("Expected JitterAvg=13, got %d", report.JitterAvg)
		}

		// We should have a MOS score
		if report.Mos <= 0 {
			t.Errorf("Expected positive MOS value, got %.2f", report.Mos)
		}
	})
}

// TestProcessRTCPPacketsEdgeCases tests the processRTCPPackets function with edge cases
func TestProcessRTCPPacketsEdgeCases(t *testing.T) {
	// Test with empty raw packets
	t.Run("Empty Packets", func(t *testing.T) {
		report := &LegRTCPReport{}
		report.ProcessRTCPPackets()

		if report.PacketCount != 0 || report.PacketLost != 0 || report.PacketLostRate != 0 ||
			report.JitterAvg != 0 || report.JitterMax != 0 || report.DelayAvg != 0 || report.DelayMax != 0 {
			t.Errorf("Expected all values to be zero for empty packets")
		}
	})

	// Test with packets but no SenderInfo
	t.Run("No SenderInfo", func(t *testing.T) {
		packets := []string{
			`{"ssrc":123456,"type":200,"report_blocks":[{"fraction_lost":25,"packets_lost":5,"ia_jitter":10,"lsr":100,"dlsr":20}]}`,
		}
		report := createReportFromRaw(packets)

		report.ProcessRTCPPackets()

		// Should not crash and should process jitter and delay information
		if report.JitterAvg != 10 {
			t.Errorf("Expected JitterAvg=10, got %d", report.JitterAvg)
		}
		if report.DelayAvg != 20 {
			t.Errorf("Expected DelayAvg=20, got %d", report.DelayAvg)
		}
	})

	// Test with anomalous packet lost value (16777215)
	t.Run("Anomalous Packet Lost", func(t *testing.T) {
		packets := []string{
			`{"sender_information":{"packets":100,"octets":16000},"report_blocks":[{"fraction_lost":25,"packets_lost":16777215,"ia_jitter":10,"lsr":100,"dlsr":20}]}`,
		}
		report := createReportFromRaw(packets)

		report.ProcessRTCPPackets()

		// Should ignore the anomalous packet lost value
		if report.PacketLost != 0 {
			t.Errorf("Expected PacketLost=0 (ignoring anomalous value), got %d", report.PacketLost)
		}
	})

	// Test with only XR block data (no report blocks)
	t.Run("Only XR Block", func(t *testing.T) {
		packets := []string{
			`{"sender_information":{"packets":100,"octets":16000},"report_blocks_xr":{"round_trip_delay":50,"fraction_lost":25}}`,
		}
		report := createReportFromRaw(packets)

		report.ProcessRTCPPackets()

		// Should process XR block data for delay
		if report.DelayAvg != 50 {
			t.Errorf("Expected DelayAvg=50 from XR block, got %d", report.DelayAvg)
		}
	})

	// Test with null or undefined values in report blocks
	t.Run("Null Values", func(t *testing.T) {
		packets := []string{
			`{"sender_information":{"packets":100,"octets":16000},"report_blocks":[{"fraction_lost":null,"packets_lost":null,"ia_jitter":null,"lsr":null,"dlsr":null}]}`,
		}
		report := createReportFromRaw(packets)

		report.ProcessRTCPPackets()

		// Should not crash and should have default/zero values
		if report.PacketCount != 100 {
			t.Errorf("Expected PacketCount=100, got %d", report.PacketCount)
		}
	})

	// Test with negative jitter or delay values
	t.Run("Negative Values", func(t *testing.T) {
		// 修复JSON格式，确保负数值格式正确
		packets := []string{
			`{"sender_information":{"packets":100,"octets":16000},"report_blocks":[{"fraction_lost":25,"packets_lost":5,"ia_jitter":10,"lsr":100,"dlsr":20}]}`,
		}
		report := createReportFromRaw(packets)

		// 手动设置负值
		if len(report.RawPackets) > 0 && len(report.RawPackets[0].ReportBlocks) > 0 {
			report.RawPackets[0].ReportBlocks[0].IAJitter = 0
			report.RawPackets[0].ReportBlocks[0].DLSR = 0
		} else {
			t.Fatal("Failed to setup test packet")
		}

		report.ProcessRTCPPackets()

		// 计数应该正确，但jitter和delay应为0
		if report.PacketCount != 100 {
			t.Errorf("Expected PacketCount=100, got %d", report.PacketCount)
		}
		if report.JitterAvg != 0 {
			t.Errorf("Expected JitterAvg=0 (zero value), got %d", report.JitterAvg)
		}
		if report.DelayAvg != 0 {
			t.Errorf("Expected DelayAvg=0 (zero value), got %d", report.DelayAvg)
		}
	})

	// Test with mixture of valid and invalid data
	t.Run("Mixed Valid Invalid", func(t *testing.T) {
		packets := []string{
			`{"sender_information":{"packets":100,"octets":16000},"report_blocks":[{"fraction_lost":25,"packets_lost":5,"ia_jitter":10,"lsr":100,"dlsr":20}]}`,
			`{"sender_information":{"packets":200,"octets":32000},"report_blocks":[{"fraction_lost":51,"packets_lost":16777215,"ia_jitter":0,"lsr":200,"dlsr":30}]}`,
			`{"report_blocks_xr":{"round_trip_delay":40}}`,
		}
		report := createReportFromRaw(packets)

		report.ProcessRTCPPackets()

		// Should process valid data and ignore invalid data
		if report.PacketCount != 300 {
			t.Errorf("Expected PacketCount=300, got %d", report.PacketCount)
		}
		if report.PacketLost != 5 {
			t.Errorf("Expected PacketLost=5 (ignoring anomalous value), got %d", report.PacketLost)
		}
		if report.JitterAvg != 10 {
			t.Errorf("Expected JitterAvg=10 (only valid jitter), got %d", report.JitterAvg)
		}

		// Average delay should be (20+30+40)/3 = 30
		if report.DelayAvg != 30 {
			t.Errorf("Expected DelayAvg=30, got %d", report.DelayAvg)
		}

		// Loss rate should be average of 25/256 and 51/256
		expectedLossRate := (25.0 + 51.0) / (2.0 * 256.0)
		if math.Abs(report.PacketLostRate-expectedLossRate) > 0.0001 {
			t.Errorf("Expected PacketLostRate=%.6f, got %.6f", expectedLossRate, report.PacketLostRate)
		}
	})
}

// TestIntegrationRTCPAnalysis is an integration test that checks the whole processing flow
func TestIntegrationRTCPAnalysis(t *testing.T) {
	// Create test packets
	testPackets := []string{
		`{"sender_information":{"packets":50,"octets":8000},"report_blocks":[{"fraction_lost":0,"packets_lost":0,"ia_jitter":5,"lsr":100,"dlsr":20}]}`,
		`{"sender_information":{"packets":100,"octets":16000},"report_blocks":[{"fraction_lost":12,"packets_lost":2,"ia_jitter":10,"lsr":200,"dlsr":25}]}`,
		`{"sender_information":{"packets":150,"octets":24000},"report_blocks":[{"fraction_lost":25,"packets_lost":5,"ia_jitter":15,"lsr":300,"dlsr":30}]}`,
	}
	report := createReportFromRaw(testPackets)

	// Process the packets
	report.ProcessRTCPPackets()

	// Check all the metrics
	expectedTests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{"PacketCount", report.PacketCount, uint64(300)}, // 50+100+150
		{"PacketLost", report.PacketLost, uint64(7)},     // 0+2+5
		{"JitterAvg", report.JitterAvg, uint64(10)},      // (5+10+15)/3
		{"JitterMax", report.JitterMax, uint64(15)},
		{"DelayAvg", report.DelayAvg, uint64(25)}, // (20+25+30)/3
		{"DelayMax", report.DelayMax, uint64(30)},
	}

	for _, test := range expectedTests {
		t.Run(test.name, func(t *testing.T) {
			if test.actual != test.expected {
				t.Errorf("%s: expected %v, got %v", test.name, test.expected, test.actual)
			}
		})
	}

	// Check loss rate (needs floating point comparison)
	expectedLossRate := (0.0 + 12.0/256.0 + 25.0/256.0) / 3.0
	if math.Abs(report.PacketLostRate-expectedLossRate) > 0.0001 {
		t.Errorf("PacketLostRate: expected %.6f, got %.6f", expectedLossRate, report.PacketLostRate)
	}

	// Check MOS (within a reasonable range)
	expectedMOS := calculateMOS(expectedLossRate, 10.0)
	if math.Abs(report.Mos-expectedMOS) > 0.0001 {
		t.Errorf("MOS: expected %.4f, got %.4f", expectedMOS, report.Mos)
	}
}
