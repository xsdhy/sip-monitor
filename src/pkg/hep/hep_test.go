package hep

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestHEP3IPv4Parsing(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected HepMsg
		wantErr  bool
	}{
		{
			name:  "Valid HEP3 IPv4 SIP message",
			input: createHEP3Packet(createBasicHEP3Header(), createIPv4Chunks(), createSIPPayloadChunk()),
			expected: HepMsg{
				Version:               3,
				IPProtocolFamily:      1,  // IPv4
				IPProtocolID:          17, // UDP
				IP4SourceAddress:      "192.168.1.10",
				IP4DestinationAddress: "192.168.1.20",
				SourcePort:            5060,
				DestinationPort:       5060,
				Timestamp:             1634567890,
				TimestampMicro:        123456,
				ProtocolType:          1, // SIP
				CaptureAgentID:        2001,
				KeepAliveTimer:        120,
				AuthenticateKey:       "testkey123",
				Body:                  []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n"),
			},
			wantErr: false,
		},
		{
			name:     "Invalid HEP3 packet (too short)",
			input:    []byte{0x48, 0x45, 0x50, 0x33},
			expected: HepMsg{},
			wantErr:  true,
		},
		{
			name:  "HEP3 packet with missing required field (IP protocol family)",
			input: createHEP3Packet(createBasicHEP3Header(), createIPv4ChunksWithoutField(0x0001), createSIPPayloadChunk()),
			expected: HepMsg{
				Version: 3,
				Body:    []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n"),
			},
			wantErr: true,
		},
		{
			name:  "HEP3 packet with scrambled chunk order",
			input: createHEP3Packet(createBasicHEP3Header(), scrambleChunks(createIPv4Chunks()), createSIPPayloadChunk()),
			expected: HepMsg{
				Version:               3,
				IPProtocolFamily:      1,
				IPProtocolID:          17,
				IP4SourceAddress:      "192.168.1.10",
				IP4DestinationAddress: "192.168.1.20",
				SourcePort:            5060,
				DestinationPort:       5060,
				Timestamp:             1634567890,
				TimestampMicro:        123456,
				ProtocolType:          1,
				CaptureAgentID:        2001,
				KeepAliveTimer:        120,
				AuthenticateKey:       "testkey123",
				Body:                  []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n"),
			},
			wantErr: false,
		},
		{
			name:  "HEP3 packet with unknown chunk type",
			input: createHEP3Packet(createBasicHEP3Header(), append(createIPv4Chunks(), createUnknownChunk()...), createSIPPayloadChunk()),
			expected: HepMsg{
				Version:               3,
				IPProtocolFamily:      1,
				IPProtocolID:          17,
				IP4SourceAddress:      "192.168.1.10",
				IP4DestinationAddress: "192.168.1.20",
				SourcePort:            5060,
				DestinationPort:       5060,
				Timestamp:             1634567890,
				TimestampMicro:        123456,
				ProtocolType:          1,
				CaptureAgentID:        2001,
				KeepAliveTimer:        120,
				AuthenticateKey:       "testkey123",
				Body:                  []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n"),
			},
			wantErr: false,
		},
		{
			name:  "HEP3 packet with incorrect chunk length",
			input: createHEP3Packet(createBasicHEP3Header(), createIPv4ChunksWithIncorrectLength(), createSIPPayloadChunk()),
			expected: HepMsg{
				Version: 3,
			},
			wantErr: true,
		},
		{
			name:  "HEP3 packet with duplicate fields",
			input: createHEP3Packet(createBasicHEP3Header(), createIPv4ChunksWithDuplicateFields(), createSIPPayloadChunk()),
			expected: HepMsg{
				Version:               3,
				IPProtocolFamily:      1,
				IPProtocolID:          17,
				IP4SourceAddress:      "192.168.1.10",
				IP4DestinationAddress: "192.168.1.20",
				SourcePort:            5060,
				DestinationPort:       5060,
				Timestamp:             1634567890,
				TimestampMicro:        123456,
				ProtocolType:          1,
				CaptureAgentID:        2001,
				KeepAliveTimer:        120,
				AuthenticateKey:       "testkey123",
				Body:                  []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n"),
			},
			wantErr: false,
		},
		{
			name:  "HEP3 packet with maximum length",
			input: createMaxLengthHEP3Packet(),
			expected: HepMsg{
				Version:               3,
				IPProtocolFamily:      1,
				IPProtocolID:          17,
				IP4SourceAddress:      "192.168.1.10",
				IP4DestinationAddress: "192.168.1.20",
				SourcePort:            5060,
				DestinationPort:       5060,
				Timestamp:             1634567890,
				TimestampMicro:        123456,
				ProtocolType:          1,
				CaptureAgentID:        2001,
				KeepAliveTimer:        120,
				AuthenticateKey:       "testkey123",
				Body:                  createLargePayload(),
			},
			wantErr: false,
		},
		{
			name:  "HEP3 packet with RTP payload",
			input: createHEP3Packet(createBasicHEP3Header(), createIPv4Chunks(), createRTPPayloadChunk()),
			expected: HepMsg{
				Version:               3,
				IPProtocolFamily:      1,
				IPProtocolID:          17,
				IP4SourceAddress:      "192.168.1.10",
				IP4DestinationAddress: "192.168.1.20",
				SourcePort:            5060,
				DestinationPort:       5060,
				Timestamp:             1634567890,
				TimestampMicro:        123456,
				ProtocolType:          5, // RTP
				CaptureAgentID:        2001,
				KeepAliveTimer:        120,
				AuthenticateKey:       "testkey123",
				Body:                  []byte{0x80, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewHepMsg(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("NewHepMsg() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil {
				compareHepMsg(t, result, &tc.expected)
			}
		})
	}
}

func compareHepMsg(t *testing.T, got, want *HepMsg) {
	t.Helper() // 将此函数标记为测试辅助函数

	// 使用反射比较结构体的每个字段
	gotVal := reflect.ValueOf(*got)
	wantVal := reflect.ValueOf(*want)
	gotType := gotVal.Type()

	for i := 0; i < gotVal.NumField(); i++ {
		fieldName := gotType.Field(i).Name
		gotField := gotVal.Field(i)
		wantField := wantVal.Field(i)

		// 特殊处理某些字段
		switch fieldName {
		case "IP4SourceAddress", "IP4DestinationAddress", "IP6SourceAddress", "IP6DestinationAddress":
			gotIP := net.IP(gotField.String())
			wantIP := net.IP(wantField.String())
			if !gotIP.Equal(wantIP) {
				t.Errorf("%s mismatch: got %v, want %v", fieldName, gotIP, wantIP)
			}
		case "Body":
			if !bytes.Equal(gotField.Bytes(), wantField.Bytes()) {
				t.Errorf("%s mismatch: got %v, want %v", fieldName, gotField.Bytes(), wantField.Bytes())
			}
		default:
			if !reflect.DeepEqual(gotField.Interface(), wantField.Interface()) {
				t.Errorf("%s mismatch: got %v, want %v", fieldName, gotField.Interface(), wantField.Interface())
			}
		}
	}
}

func createBasicHEP3Header() []byte {
	header := make([]byte, 6)
	header[0] = 0x48 // 'H'
	header[1] = 0x45 // 'E'
	header[2] = 0x50 // 'P'
	header[3] = 0x33 // '3'
	// 总长度先设为0，后面再更新
	return header
}

func createIPv4Chunks() []byte {
	chunks := []byte{}

	// IP协议族 (IPv4)
	chunks = append(chunks, createChunk(0x0001, []byte{0x01})...)

	// IP协议ID (UDP)
	chunks = append(chunks, createChunk(0x0002, []byte{17})...)

	// IPv4源地址
	chunks = append(chunks, createChunk(0x0003, net.ParseIP("192.168.1.10").To4())...)

	// IPv4目标地址
	chunks = append(chunks, createChunk(0x0004, net.ParseIP("192.168.1.20").To4())...)

	// 源端口
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, 5060)
	chunks = append(chunks, createChunk(0x0007, portBytes)...)

	// 目标端口
	chunks = append(chunks, createChunk(0x0008, portBytes)...)

	// 时间戳
	tsBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tsBytes, 1634567890)
	chunks = append(chunks, createChunk(0x0009, tsBytes)...)

	// 微秒时间戳
	tsMicroBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tsMicroBytes, 123456)
	chunks = append(chunks, createChunk(0x000a, tsMicroBytes)...)

	// 协议类型 (SIP)
	chunks = append(chunks, createChunk(0x000b, []byte{0x01})...)

	// 捕获代理ID
	captureAgentIDBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(captureAgentIDBytes, 2001)
	chunks = append(chunks, createChunk(0x000c, captureAgentIDBytes)...)

	// 保活定时器
	keepAliveTimerBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(keepAliveTimerBytes, 120)
	chunks = append(chunks, createChunk(0x000d, keepAliveTimerBytes)...)

	// 认证密钥
	chunks = append(chunks, createChunk(0x000e, []byte("testkey123"))...)

	return chunks
}

func createSIPPayloadChunk() []byte {
	payload := []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n")
	return createChunk(0x000f, payload)
}

func createChunk(chunkType uint16, data []byte) []byte {
	chunk := make([]byte, 6+len(data))
	binary.BigEndian.PutUint16(chunk[0:2], 0x0000) // Vendor ID (0 for IETF)
	binary.BigEndian.PutUint16(chunk[2:4], chunkType)
	binary.BigEndian.PutUint16(chunk[4:6], uint16(6+len(data)))
	copy(chunk[6:], data)
	return chunk
}

func createHEP3Packet(header, chunks, payload []byte) []byte {
	packet := append(header, chunks...)
	packet = append(packet, payload...)
	binary.BigEndian.PutUint16(packet[4:6], uint16(len(packet)))
	return packet
}

func createIPv4ChunksWithoutField(fieldToRemove uint16) []byte {
	allChunks := createIPv4Chunks()
	var result []byte

	for i := 0; i < len(allChunks); {
		if i+6 > len(allChunks) {
			// 如果剩余的字节不足以构成一个完整的 chunk 头部，就退出循环
			break
		}

		chunkType := binary.BigEndian.Uint16(allChunks[i+2 : i+4])
		chunkLength := binary.BigEndian.Uint16(allChunks[i+4 : i+6])

		if int(chunkLength) > len(allChunks[i:]) {
			// 如果 chunk 的长度超出了剩余的字节数，就退出循环
			break
		}

		if chunkType != fieldToRemove {
			result = append(result, allChunks[i:i+int(chunkLength)]...)
		}

		i += int(chunkLength)
	}

	return result
}

func scrambleChunks(chunks []byte) []byte {
	// 简单的打乱算法，实际使用时可以使用更复杂的随机算法
	result := make([]byte, len(chunks))
	copy(result, chunks)
	for i := len(result) - 1; i > 0; i-- {
		j := i / 2
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func createUnknownChunk() []byte {
	return createChunk(0xFFFF, []byte{0x00, 0x01, 0x02, 0x03})
}

func createIPv4ChunksWithIncorrectLength() []byte {
	chunks := createIPv4Chunks()
	// 修改第一个chunk的长度字段
	binary.BigEndian.PutUint16(chunks[4:6], 1000)
	return chunks
}

func createIPv4ChunksWithDuplicateFields() []byte {
	chunks := createIPv4Chunks()
	// 添加一个重复的源IP地址字段
	duplicateChunk := createChunk(0x0003, net.ParseIP("192.168.1.11").To4())
	return append(chunks, duplicateChunk...)
}

func createMaxLengthHEP3Packet() []byte {
	header := createBasicHEP3Header()
	chunks := createIPv4Chunks()
	payload := createLargePayload()
	return createHEP3Packet(header, chunks, createChunk(0x000f, payload))
}

func createLargePayload() []byte {
	payload := make([]byte, 65000) // 假设最大长度为65000字节
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	return payload
}

func createRTPPayloadChunk() []byte {
	rtpPacket := []byte{
		0x80, 0x00, 0x00, 0x01, // RTP header
		0x00, 0x00, 0x00, 0x00, // Timestamp
		0x00, 0x00, 0x00, 0x00, // SSRC
	}
	return createChunk(0x000f, rtpPacket)
}
