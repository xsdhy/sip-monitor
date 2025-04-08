package util

import (
	"testing"
)

func TestSanitizeRawText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有ASCII可打印字符",
			input:    "Hello, World! 123",
			expected: "Hello, World! 123",
		},
		{
			name:     "包含基本控制字符",
			input:    "Hello\r\nWorld\tTest",
			expected: "Hello\r\nWorld\tTest",
		},
		{
			name:     "包含中文字符",
			input:    "Hello世界",
			expected: "Hello[0x4E16][0x754C]",
		},
		{
			name:     "混合字符",
			input:    "ABC123!@#\r\n\t",
			expected: "ABC123!@#\r\n\t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeRawText(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeRawText() = %v, want %v", result, tt.expected)
			}
		})
	}

	// 测试替换字符处理
	t.Run("包含不可见字符", func(t *testing.T) {
		// 使用byte数组构造带有不可见字符的字符串
		rawBytes := []byte{'H', 'e', 'l', 'l', 'o', 0x81, 'W', 'o', 'r', 'l', 'd'}
		input := string(rawBytes)

		// Go自动将无效UTF-8序列替换为Unicode替换字符(U+FFFD)
		expected := "Hello[0xFFFD]World"
		result := SanitizeRawText(input)

		if result != expected {
			t.Errorf("SanitizeRawText() = %v, want %v", result, expected)
		}
	})

	t.Run("包含多种不可见字符", func(t *testing.T) {
		// 构造测试样例中的"测试\x81~\x06\x1DSI"
		rawBytes := []byte{0xE6, 0xB5, 0x8B, 0xE8, 0xAF, 0x95, 0x81, '~', 0x06, 0x1D, 'S', 'I'}
		input := string(rawBytes)

		// Go自动将无效UTF-8序列替换为Unicode替换字符(U+FFFD)
		expected := "[0x6D4B][0x8BD5][0xFFFD]~[0x6][0x1D]SI"
		result := SanitizeRawText(input)

		if result != expected {
			t.Errorf("SanitizeRawText() = %v, want %v", result, expected)
		}
	})
}
