package arrkit

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

// 从 string 数组中删除一个指定值，返回删除后的新数组
func ArrayStringDeleteVal(val string, array []string) []string {
	// 找出值的 index
	index := 0
	count := len(array)
	for i := 0; i < count; i++ {
		if val == array[i] {
			index = i
			break
		}
	}
	// 删除
	array = append(array[:index], array[index+1:]...)
	return array
}

// 从数组中删除一个元素
// index 要删除的元素下标
func ArrayStringDeleteOne(index int, array []string) []string {
	array = append(array[:index], array[index+1:]...)
	return array
}

// 判断符合数据类型中是否存在某个值
// 支持的类型: slice 、array 、map
func InArray(value interface{}, array interface{}) bool {
	val := reflect.ValueOf(array)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(value, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(value, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: 只支持类型为 slice、array、map 的数据验证")
	}
	return false
}

// Int 数组倒序
func IntArrayDesc(arr []int) (ret []int) {
	sort.Ints(arr)
	for i := len(arr) - 1; i >= 0; i-- {
		ret = append(ret, arr[i])
	}
	return ret
}

// 数组转逗号拼接
func JoinWithCommas(numbers []int) string {
	var strNumbers []string
	for _, num := range numbers {
		strNumbers = append(strNumbers, fmt.Sprint(num))
	}
	return strings.Join(strNumbers, ",")
}

// 用给定的键值填充
func Array_fill(startIndex int, num uint, value interface{}) map[int]interface{} {
	m := make(map[int]interface{})
	var i uint
	for i = 0; i < num; i++ {
		m[startIndex] = value
		startIndex++
	}
	return m
}

// 反转/交换数组中所有的键名以及它们关联的键值。
func Array_flip(m map[interface{}]interface{}) map[interface{}]interface{} {
	n := make(map[interface{}]interface{})
	for i, v := range m {
		n[v] = i
	}
	return n
}

// 返回包含数组中所有键名的一个新数组
func Array_keys(elements map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(elements))
	for key := range elements {
		keys[i] = key
		i++
	}
	return keys
}

// 以数值下标重建索引
func Array_values(elements map[interface{}]interface{}) []interface{} {
	i, vals := 0, make([]interface{}, len(elements))
	for _, val := range elements {
		vals[i] = val
		i++
	}
	return vals
}

// 数组合并
func Array_merge(ss ...[]interface{}) []interface{} {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]interface{}, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// 把一个数组分割为新的数组块
func Array_chunk(s []interface{}, size int) [][]interface{} {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

// 数组填充
func Array_pad(s []interface{}, size int, val interface{}) []interface{} {
	if size == 0 || (size > 0 && size < len(s)) || (size < 0 && size > -len(s)) {
		return s
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(s)
	tmp := make([]interface{}, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		return append(s, tmp...)
	}
	return append(tmp, s...)
}

// 在数组中根据条件切出一段值
func Array_slice(s []interface{}, offset, length uint) []interface{} {
	if offset > uint(len(s)) {
		panic("offset: the offset is less than the length of s")
	}
	end := offset + length
	if end < uint(len(s)) {
		return s[offset:end]
	}
	return s[offset:]
}

// 随机返回数组中的一个键名
func Array_rand(elements []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]interface{}, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// 从数组中取出所有某个键名的值
// 支持数据类型：
//
//	[]map[string]interface {}
//	map[string]map[string]interface {}
func Array_column(array interface{}, columnKey string) interface{} {
	theType := reflect.TypeOf(array).String()
	var returnColumns interface{}
	switch theType {
	case "[]map[string]interface {}":
		input := array.([]map[string]interface{})
		columns := make([]interface{}, 0, len(input))
		for _, val := range input {
			if v, ok := val[columnKey]; ok {
				columns = append(columns, v)
			}
		}
		returnColumns = columns
	case "map[string]map[string]interface {}":
		input := array.(map[string]map[string]interface{})
		columns := make([]interface{}, 0, len(input))
		for _, val := range input {
			if v, ok := val[columnKey]; ok {
				columns = append(columns, v)
			}
		}
		returnColumns = columns
	}
	return returnColumns
}

// 数组尾部压入栈
func Array_push(s *[]interface{}, elements ...interface{}) int {
	*s = append(*s, elements...)
	return len(*s)
}

// 删除数组中的最后一个元素，出栈
func Array_pop(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	ep := len(*s) - 1
	e := (*s)[ep]
	*s = (*s)[:ep]
	return e
}

// 数组头部插入一个元素
func Array_unshift(s *[]interface{}, elements ...interface{}) int {
	*s = append(elements, *s...)
	return len(*s)
}

// 数组头部删除一个元素
func Array_shift(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// 通过合并两个数组来创建一个新数组，其中的一个数组元素为键名，另一个数组的元素为键值。
func Array_combine(s1, s2 []interface{}) map[interface{}]interface{} {
	if len(s1) != len(s2) {
		panic("the number of elements for each slice isn'taa equal")
	}
	m := make(map[interface{}]interface{}, len(s1))
	for i, v := range s1 {
		m[v] = s2[i]
	}
	return m
}

// 数组顺序反转，以相反的元素顺序返回数组
func Array_reverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// 数组切割为字符串
func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// 检查数组/切片/map中是否存在指定的键名
func Array_key_exists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]
	return ok
}
