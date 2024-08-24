package mathkit

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 取随机数
// 范围: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// 保留指定位随机小数
func RandDecimals(min, max float64, precision ...int) float64 {
	prec := 2
	if len(precision) > 0 {
		prec = precision[0]
	}
	result := min + rand.Float64()*(max-min)
	scale := math.Pow(10, float64(prec))
	return math.Round(result*scale) / scale
}

// 对浮点数进保留几位小数
func Floor(value float64, decimal int) float64 {
	decimalStr := strconv.Itoa(decimal)
	finalValue, _ := strconv.ParseFloat(fmt.Sprintf("%."+decimalStr+"f", value), 64)
	return finalValue
}

// 对浮点数进保留几位小数，0 填充小数位，0 填充只能是 string 类型才能显示完整
func FloorWithZeroPad(value float64, decimal int) string {
	decimalStr := strconv.Itoa(decimal)
	finalValue, _ := strconv.ParseFloat(fmt.Sprintf("%."+decimalStr+"f", value), 64)
	// 小数不足 decimal 位长度时 0 填充
	return fmt.Sprintf("%.4f", finalValue)
}

//////////// 数学计算相关函数 ////////////

// 取绝对值
func Abs(number float64) float64 {
	return math.Abs(number)
}

// 对浮点数进行四舍五入
// args[0] 保留几位小数
func Round(value float64, args ...int) float64 {
	if len(args) > 0 {
		value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
		return value
	}
	return math.Floor(value + 0.5)
}

// 进1取整
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// 返回几个指定值中的最大值
func Max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}
	return max
}

// 返回几个指定值中的最小值
func Min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min
}

// 把十进制数转换为二进制数
func Decbin(number int64) string {
	return strconv.FormatInt(number, 2)
}

// 把二进制转换为十进制
func Bindec(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 10), nil
}

// 把十六进制值的字符串转换为 ASCII 字符
func Hex2bin(data string) (string, error) {
	i, err := strconv.ParseInt(data, 16, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 2), nil
}

// 把 ASCII 字符的字符串转换为十六进制值
func Bin2hex(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 16), nil
}

// 把十进制数转换为十六进制数
func Dechex(number int64) string {
	return strconv.FormatInt(number, 16)
}

// 把十六进制转换为十进制
func Hexdec(str string) (int64, error) {
	return strconv.ParseInt(str, 16, 0)
}

// 把十进制转换为八进制
func Decoct(number int64) string {
	return strconv.FormatInt(number, 8)
}

// 八进制转十进制
func Octdec(str string) (int64, error) {
	return strconv.ParseInt(str, 8, 0)
}

// 把十六进制数转换为八进制数
func Base_convert(number string, frombase, tobase int) (string, error) {
	i, err := strconv.ParseInt(number, frombase, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, tobase), nil
}

// 基于微秒计的当前时间生成唯一ID
func Uniqid(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%08x%05x", prefix, now.Unix(), now.UnixNano()%0x100000)
}

// 通过千位分组来格式化数字。
// decimals: 设置保留几位小数
// decPoint: 设置小数点的分隔符。
// thousandsSep: 设置千位分隔符。
func Number_format(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}
