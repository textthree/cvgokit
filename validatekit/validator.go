package validatekit

import "regexp"

// 判断给定的字符串是否符合邮箱格式
func IsEmail(email string) bool {
	// 定义邮箱的正则表达式
	// 这个正则表达式是简化的，用于演示。在实际应用中，你可能需要更复杂的正则表达式来确保更全面的验证。
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
