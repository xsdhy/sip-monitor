package main

import (
	"fmt"
	"sip-monitor/src/pkg/util"
)

func main() {
	// 测试几个含有不支持字符的字符串
	testStrings := []string{
		"Hello, World!",        // 纯ASCII
		"Hello\r\nWorld\tTest", // 含控制字符
		"Hello\x81World",       // 含不支持字符
		"Hello世界",              // 含中文
		"测试\x81~\x06\x1DSI",    // 混合非ASCII字符
		"Incorrect string value: '\\x81~\\x06\\x1DSI...'", // 类似错误信息中的字符
	}

	// 测试每个字符串并打印结果
	for i, str := range testStrings {
		sanitized := util.SanitizeRawText(str)
		fmt.Printf("测试 %d:\n原文: %s\n处理后: %s\n\n", i+1, str, sanitized)
	}

	// 测试一个模拟的SIP消息
	rawSip := `INVITE sip:1001@192.168.1.1 SIP/2.0
Via: SIP/2.0/UDP 192.168.1.2:5060;branch=z9hG4bK-524287-1---a511d8459a80a082
From: <sip:1002@192.168.1.2>;tag=a511d8459a80a082
To: <sip:1001@192.168.1.1>
Call-ID: AEY9a6fafb7a4bbdd3c97e555
CSeq: 1 INVITE
Contact: <sip:1002@192.168.1.2:5060>
Content-Type: application/sdp
Content-Length: 138

` + string([]byte{118, 61, 48, 13, 10, 111, 61, 45, 32, 48, 32, 48, 32, 73, 78, 32, 73, 80, 52, 32, 49, 57, 50, 46, 49, 54, 56, 46, 49, 46, 50, 13, 10, 115, 61, 45, 13, 10, 99, 61, 73, 78, 32, 73, 80, 52, 32, 49, 57, 50, 46, 49, 54, 56, 46, 49, 46, 50, 13, 10, 116, 61, 48, 32, 48, 13, 10, 109, 61, 97, 117, 100, 105, 111, 32, 51, 48, 50, 48, 32, 82, 84, 80, 47, 65, 86, 80, 32, 48, 13, 10, 97, 61, 114, 116, 112, 109, 97, 112, 58, 48, 32, 80, 67, 77, 85, 47, 56, 48, 48, 48, 13, 10, 0})

	sanitizedSip := util.SanitizeRawText(rawSip)
	fmt.Println("SIP消息测试:")
	fmt.Printf("处理后: %s\n", sanitizedSip)
}
