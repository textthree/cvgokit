package gokit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 获取模块名称
func GetModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			// 提取 module 名称
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("获取模块名称出错，当前目录找不到 go.mod 文件")
}
