package bizkit

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成唯一订单号
func GenerateOrderSn(prefix ...string) string {
	if len(prefix) == 0 {
		prefix = append(prefix, "SN")
	}
	now := time.Now()
	// Format time as YYYYMMDDHHMMSS
	timePart := now.Format("20060102150405")
	// Generate a random number with 6 digits
	randomPart := rand.Intn(1000000)
	randomPartStr := fmt.Sprintf("%06d", randomPart)
	// Combine prefix, time part and random part
	orderID := fmt.Sprintf("%s%s%s", prefix[0], timePart, randomPartStr)
	return orderID
}
