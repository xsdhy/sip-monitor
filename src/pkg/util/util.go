package util

import (
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const DayFormat = "2006_01_02"

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

func GetDay(days int) string {
	return time.Now().AddDate(0, 0, days).Format(DayFormat)
}

// ParseInt64 将字符串转换为int64类型
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// HashPassword 使用bcrypt对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码是否与哈希值匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
