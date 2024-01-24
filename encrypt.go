package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// 签名算法
func Sign(params map[string]string, secret string) string {
	// 对参数按照参数名进行升序排列
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 将参数和值进行拼接，用"&"连接
	var s []string
	for _, k := range keys {
		v := params[k]
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}
	signStr := strings.Join(s, "&")

	// URL编码
	signStr = url.QueryEscape(signStr)

	// 拼接密钥，并对字符串进行哈希
	h := md5.New()
	h.Write([]byte(secret + signStr))
	signBytes := h.Sum(nil)

	// 将哈希值转换为16进制字符串
	return hex.EncodeToString(signBytes)
}
