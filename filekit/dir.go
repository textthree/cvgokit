package filekit

import (
	"fmt"
	"os"
)

// 创建目录，如果目录不存在则创建，存在不创建
func MkDir(dir string, modeArg ...os.FileMode) {
	mode := os.ModePerm
	if len(modeArg) > 0 {
		mode = modeArg[0]
	}
	var err error
	if _, err = os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, mode); err != nil {
			fmt.Println("创建目录失败：" + dir)
		}
	}
}

// EnsureDirExists 确保目标目录存在，如果不存在则创建它
func EnsureDirExists(dir string) error {
	// 检查目录是否存在，不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to add directory: %v", err)
		}
	}
	return nil
}

// 检查目录是否存在
func DirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// 删除给定的文件或目录，如果不存在则直接返回 nil，如果存在则删除
func DeleteDirOrFile(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 路径不存在，直接返回
		return nil
	} else if err != nil {
		// 其他错误，如权限问题
		return fmt.Errorf("error checking if path exists: %v", err)
	}
	// 路径存在，执行删除操作
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("error deleting path: %v", err)
	}
	return nil
}
