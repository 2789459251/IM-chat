package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写
func Md5Encode(date string) string {
	h := md5.New()
	h.Write([]byte(date))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)

}

// 大写
func MD5Encode(date string) string {
	return strings.ToUpper(Md5Encode(date))
}

// 加密
func MakePasswd(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, passwd string) bool {
	Md5Encode(plainpwd + salt)
	return Md5Encode(plainpwd+salt) == passwd
}
