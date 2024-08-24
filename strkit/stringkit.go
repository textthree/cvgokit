package strkit

import (
	"bytes"
	"encoding/json"
	"html"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// 类型转字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Tostring(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

// 中横线拼接的字符串转驼峰式
func MethodNameToCamel(str string) string {
	// 切成数组
	stringArray := strings.Split(str, "-")
	// 首字母转大写
	for key, value := range stringArray {
		stringArray[key] = Ucfirst(value)
	}
	// 拼接
	returnString := strings.Join(stringArray, "")
	return returnString
}

// 将驼峰 PascalCase 或 CamelCase 转换为 snake_case
func CamelToSnake(s string) string {
	// 定义正则表达式，用于匹配大写字母和缩写
	var matchFirstCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	var matchAllCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")

	// 添加下划线并转换为小写
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)

	return snake
}

// 下划线转大驼峰
// SnakeToPascalCase 将 snake_case 转换为 PascalCase
func SnakeToPascalCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		if len(word) > 0 {
			// 将第一个字母转换为大写，其他保持不变
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, "")
}

// SnakeToCamelCase 将 snake_case 转换为 CamelCase
func SnakeToCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		if len(word) > 0 {
			if i == 0 {
				// 第一个单词首字母小写
				words[i] = strings.ToLower(string(word[0])) + word[1:]
			} else {
				// 其他单词首字母大写
				words[i] = strings.ToUpper(string(word[0])) + word[1:]
			}
		}
	}
	return strings.Join(words, "")
}

// 生成唯一字符串，可以传入参数将唯一数值生成唯一字符串，比如手机号
// 用固定数值数值传按固定规则混淆，生成唯一字符串
func UniqueString(args ...string) string {
	var str string
	if len(args) == 0 {
		str = string(int(time.Now().Unix()))
	} else {
		str = args[0]
	}
	uniqueStr := Strtr(str, "1234567890", "huvtwkaemx")
	// 把顺序颠倒过来，然后再在中间找两个位置位置各插入一个字母
	var buffer bytes.Buffer
	for i := len(uniqueStr) - 1; i >= 0; i-- {
		buffer.WriteString(string(uniqueStr[i]))
		if i == 9 {
			buffer.WriteString(string(uniqueStr[3]))
		}
		if i == 5 {
			buffer.WriteString(string(uniqueStr[10]))
		}
	}
	uniqueStr = buffer.String()
	// 再交换一下顺序,藏头藏尾
	temp := Explode("", uniqueStr)
	_temp := temp[1]
	temp[1] = temp[7]
	temp[7] = _temp

	_temp = temp[2]
	temp[2] = temp[9]
	temp[9] = _temp

	_temp = temp[11]
	temp[11] = temp[5]
	temp[5] = _temp

	_temp = temp[12]
	temp[3] = temp[12]
	temp[12] = _temp

	var buffer2 bytes.Buffer
	for _, v := range temp {
		buffer2.WriteString(v)
	}
	uniqueStr = buffer.String()
	// 最后拼接个前缀
	uniqueStr = "mt" + uniqueStr
	return uniqueStr
}

// 生成唯一数字串，把时间戳的第一位砍掉，换成0-9的随机，第一位是1，到2033年，时间戳第一位变成2
func UniqueNumber() string {
	int64 := time.Now().Unix()
	timestamp := strconv.FormatInt(int64, 10)
	behead := timestamp[1:len(timestamp)]
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 随机种子
	uniqueNumber := strconv.FormatInt(r.Int63n(8)+1, 10) + behead
	// 交换一下位置
	temp := Explode("", uniqueNumber)
	_temp := temp[1]
	temp[1] = temp[9]
	temp[9] = _temp

	_temp = temp[2]
	temp[2] = temp[7]
	temp[7] = _temp

	_temp = temp[4]
	temp[4] = temp[6]
	temp[6] = _temp

	var buffer bytes.Buffer
	for _, v := range temp {
		buffer.WriteString(v)
	}
	uniqueNumber = buffer.String()
	return uniqueNumber
}

// 生成随机字符串
func CreateNonceStr(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 获取文件后缀
func GetSuffix(fileName string) string {
	s := Explode(".", fileName)
	last := s[len(s)-1]
	if last == fileName {
		return ""
	}
	return last
}

/**
 * 去除字符串首/尾逗号
 * @param str 要操作的字符串
 * @param mode 可选值：left 去除左边逗号，right 去除右边逗号。默认去除左右两边
 */
func TrimComma(str string, arg ...string) string {
	if str == "" || str == "," {
		return ""
	}
	mode := "ALL"
	if len(arg) > 0 {
		mode = arg[0]
	}
	str = Trim(str)
	ret := str
	switch strings.ToUpper(mode) {
	case "LEFT":
		if Substr(str, 0, 1) == "," {
			ret = str[1:]
		}
	case "RIGHT":
		if str[len(str)-1:] == "," {
			ret = str[0 : len(str)-1]
		}
	case "ALL":
		if Substr(str, 0, 1) == "," {
			ret = str[1:]
		}
		if ret[len(ret)-1:] == "," {
			ret = ret[0 : len(ret)-1]
		}
	}
	if ret == "" {
		ret = str
	}
	return ret
}

// 字符串切割成数组，返回字符串切片
func Explode(delimiter, str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, delimiter)
}

// 搜索字符串在另一字符串中是否存在，如果存在则返回该字符串及剩余部分，否则返回 FALSE。
func Strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

// 字符串翻译函数，转换字符串中特定的字符。
// 如果params ...interface{}只传一个参数，类型是: map[string]string
// 例如：Strtr("baab", map[string]string{"ab": "01"}) 返回 "ba01"
// 如果params ...interface{}传两个参数, 类型是：string, string
// Strtr("baab", "ab", "01") 返回 "1001", a => 0; b => 1
func Strtr(haystack string, params ...interface{}) string {
	ac := len(params)
	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)
		if length == 0 {
			return haystack
		}
		oldnew := make([]string, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}
		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)
		if trlen > lt {
			trlen = lt
		}
		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8
		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}
			return string(str)
		}
		// trlen != 1
		for {
			xlat[j] = j
			if j++; j == 0 {
				break
			}
		}
		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}
		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}
		return string(str)
	}

	return haystack
}

// 查找字符串在另一字符串中首次出现的位置（区分大小写）
// 在 str 中查找 find
func Strpos(str, find string, offsetArg ...int) int {
	var offset int
	if len(offsetArg) == 0 {
		offset = 0
	} else {
		offset = offsetArg[0]
	}
	length := len(str)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}
	pos := strings.Index(str[offset:], find)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// 查找字符串在另一字符串中首次出现的位置（不区分大小写）
func Stripos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	haystack = haystack[offset:]
	if offset < 0 {
		offset += length
	}
	pos := strings.Index(strings.ToLower(haystack), strings.ToLower(needle))
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// 查找字符串在另一字符串中最后一次出现的位置（区分大小写）
// haystack : 被查找的字符串
// needle : 要在haystack中查找的字符串
// args[0] :  可选，规定从何处开始搜索
func Strrpos(haystack, needle string, args ...int) int {
	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos = strings.LastIndex(haystack, needle)
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// 查找字符串在另一字符串中最后一次出现的位置（不区分大小写）
func Strripos(haystack, needle string, offset int) int {
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos = strings.LastIndex(strings.ToLower(haystack), strings.ToLower(needle))
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// 字符串替换，在 subject 中将 old 替换成 new
func StrReplace(old, new, subject string, count ...int) string {
	num := -1
	if len(count) > 0 {
		num = count[0]
	}
	// -1 代表替换全部，0 代表不做替换，1 代表只替换一次
	return strings.Replace(subject, old, new, num)
}

// 字符串转大写
func Strtoupper(str string) string {
	return strings.ToUpper(str)
}

// 字符串转小写
func Strtolower(str string) string {
	return strings.ToLower(str)
}

// 字符串首字母转化为大写
func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

// 首字母转小写
func Lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

// 单词首字母转大写
func Ucwords(str string) string {
	return strings.Title(str)
}

// 字符串截取
func Substr(str string, start uint, length int) string {
	if start < 0 || length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

// 字符串反转
func Strrev(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 去除字符串两边空格
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, characterMask[0])
}

// 去除字符串左边空格
func Ltrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	return strings.TrimLeft(str, characterMask[0])
}

// 去除字符串右边空格
func Rtrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	return strings.TrimRight(str, characterMask[0])
}

func ExplodeAndTrim(delimiter, str string) []string {
	arr := strings.Split(str, delimiter)
	var ret []string
	for _, v := range arr {
		ret = append(ret, Trim(v))
	}
	return ret
}

// string 转 int
func ParseInt(str string) int {
	data, _ := strconv.Atoi(str)
	return data
}

// string 转 int8
func ParseInt8(str string) int8 {
	data, _ := strconv.ParseInt(str, 10, 8)
	return int8(data)
}

// string 转 int32
func ParseInt32(str string) int32 {
	data, _ := strconv.ParseInt(str, 10, 32)
	return int32(data)
}

// string 转 int64
func ParseInt64(str string) int64 {
	// 如果是小数，
	data, _ := strconv.ParseInt(str, 10, 64)
	return data
}

func StringToFloat64(str string) float64 {
	data, _ := strconv.ParseFloat(str, 64)
	return data
}

// 判断字符串是否已 xxx 开头
func StartWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// 判断字符串是否已 xxx 结尾
func EndWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// 把字符串重复指定次数
func StrRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// 去除字符串中的所有空白字符串，包括空格 \n \r \t
func RemoveSpace(input string) string {
	result := strings.ReplaceAll(input, " ", "")
	result = strings.ReplaceAll(result, "\t", "")
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\r", "")
	return result
}

// 把字符串按指定长度分隔并拼接分隔符
// 例如在每个字符后分割一次字符串，并在每个分割后添加 "."：
// str := "tangsan";
// Chunk_split(str,1,".");  得到：taa.a.n.g.s.a.n.
func Chunk_split(body string, chunklen uint, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, erunes := []rune(body), []rune(end)
	l := uint(len(runes))
	if l <= 1 || l < chunklen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(erunes))
	var i uint
	for i = 0; i < l; i += chunklen {
		if i+chunklen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunklen]...)
		}
		ns = append(ns, erunes...)
	}
	return string(ns)
}

// 按指定长度对字符串进行拆分，
// 用于换行，例如： Wordwrap(str,15,"<br>\n");
func Wordwrap(str string, width uint, br string, cut bool) string {
	strlen := len(str)
	brlen := len(br)
	linelen := int(width)

	if strlen == 0 {
		return ""
	}
	if brlen == 0 {
		panic("break string cannot be empty")
	}
	if linelen == 0 && cut {
		panic("can'taa force cut when width is zero")
	}

	current, laststart, lastspace := 0, 0, 0
	var ns []byte
	for current = 0; current < strlen; current++ {
		if str[current] == br[0] && current+brlen < strlen && str[current:current+brlen] == br {
			ns = append(ns, str[laststart:current+brlen]...)
			current += brlen - 1
			lastspace = current + 1
			laststart = lastspace
		} else if str[current] == ' ' {
			if current-laststart >= linelen {
				ns = append(ns, str[laststart:current]...)
				ns = append(ns, br[:]...)
				laststart = current + 1
			}
			lastspace = current
		} else if current-laststart >= linelen && cut && laststart >= lastspace {
			ns = append(ns, str[laststart:current]...)
			ns = append(ns, br[:]...)
			laststart = current
			lastspace = current
		} else if current-laststart >= linelen && laststart < lastspace {
			ns = append(ns, str[laststart:lastspace]...)
			ns = append(ns, br[:]...)
			lastspace++
			laststart = lastspace
		}
	}

	if laststart != current {
		ns = append(ns, str[laststart:current]...)
	}
	return string(ns)
}

// 获取中文字符串长度
func Mb_strlen(str string) int {
	return utf8.RuneCountInString(str)
}

// 把字符串重复指定次数
func Str_repeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// 随机地打乱字符串中的所有字符：
func Str_shuffle(str string) string {
	runes := []rune(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// 将ASCII码值转化为字符串。
// 此函数与PHP的Mb_chr()函数转换结果一致，与php的chr()转换结果不一致
// 因为golang统一是utf-8编码，rune uses UTF-8，ASCII码值在127以下，127一下是可以和php对等，超过127的ASCII值翻译就无法对等了
// 曹尼玛在移植php discuz经典加解密算法过来时，php和golang的这个差别折腾了半天
func Chr(ascii int) string {
	return string(ascii)
}

// 将字符串转化为ASCII码值
func Ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

// 将换行符转成HTML的<br/>标签
// \n\r, \r\n, \r, \n
func Nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

// 转义引号
func Addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// 删除由 addslashes() 函数添加的反斜杠。
func Stripslashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// 在字符串中某些预定义的字符前添加反斜杠。
func Quotemeta(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// 把字符转换为 HTML 实体
func Htmlentities(str string) string {
	return html.EscapeString(str)
}

// 把Htmlentities()处理成的实体转回成字符串
func HTMLEntityDecode(str string) string {
	return html.UnescapeString(str)
}
