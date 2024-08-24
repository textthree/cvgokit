package filekit

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// 获取文件信息
// filename 文件路径
func Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

// 获取文件路径相关信息
// -1: all; 1: dirname; 2: basename; 4: extension; 8: filename
// 用法:
// Pathinfo("/home/go/path/src/php2go/php2go.go", 1|2|4|8)
func Pathinfo(path string, options int) map[string]string {
	if options == -1 {
		options = 1 | 2 | 4 | 8
	}
	info := make(map[string]string)
	if (options & 1) == 1 {
		info["dirname"] = filepath.Dir(path)
	}
	if (options & 2) == 2 {
		info["basename"] = filepath.Base(path)
	}
	if ((options & 4) == 4) || ((options & 8) == 8) {
		basename := ""
		if (options & 2) == 2 {
			basename, _ = info["basename"]
		} else {
			basename = filepath.Base(path)
		}
		p := strings.LastIndex(basename, ".")
		filename, extension := "", ""
		if p > 0 {
			filename, extension = basename[:p], basename[p+1:]
		} else if p == -1 {
			filename = basename
		} else if p == 0 {
			extension = basename[p+1:]
		}
		if (options & 4) == 4 {
			info["extension"] = extension
		}
		if (options & 8) == 8 {
			info["filename"] = filename
		}
	}
	return info
}

// 判断文件或目录是否存在
func File_exists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 判断是否是一个目录
func Is_dir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

// 获取文件大小
func Filesize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}

// copy一个文件
func Copy(source, dest string) (bool, error) {
	fd1, err := os.Open(source)
	if err != nil {
		return false, err
	}
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	_, e := io.Copy(fd2, fd1)
	if e != nil {
		return false, e
	}
	return true, nil
}

// 判断文件是否可读
func IsReadable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return false
	}
	return true
}

// 判断文件是否可写
func IsWriteable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_WRONLY, 0)
	if err != nil {
		return false
	}
	return true
}

// 创建一个文件
func Touch(filename string) (bool, error) {
	_, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 返回文件绝对路径
func Realpath(path string) (string, error) {
	return filepath.Abs(path)
}

// 返回路径中的文件名部分
func Basename(path string) string {
	return filepath.Base(path)
}

// 更改权限
func Chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// 从文件指针中读入一行并解析 CSV 字段
func Fgetcsv(handle *os.File, length int, delimiter rune) ([][]string, error) {
	reader := csv.NewReader(handle)
	reader.Comma = delimiter
	// TODO length limit
	return reader.ReadAll()
}

// 返回一个包含匹配指定模式的文件名或目录的数组
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
