package filekit

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetParentDir 获取级目录的绝对路径
// deep 指定上几级
func GetParentDir(deep ...int8) string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var count int8 = 1
	if len(deep) > 0 {
		count = deep[0]
	}
	targetDir := currentDir
	for i := int8(0); i < count; i++ {
		targetDir = filepath.Dir(targetDir)
	}
	return targetDir
}
