package filekit

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipDirectory 将 srcDir 目录打包压缩，并输出到 dstZipFile
// e.g. ZipDirectory(dist, dist+"dist.zip")
func ZipDirectory(srcDir, dstZipFile string) error {
	// 创建目标 ZIP 文件
	zipFile, err := os.Create(dstZipFile)
	if err != nil {
		return fmt.Errorf("failed to add zip file: %v", err)
	}
	defer zipFile.Close()

	// 创建 ZIP Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录并将文件添加到 ZIP
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是目标 ZIP 文件，避免死循环
		if path == dstZipFile {
			return nil
		}

		// 获取相对路径并将其标准化为 Unix 路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		// 如果是目录，添加到 ZIP 但不压缩内容
		if info.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			if err != nil {
				return err
			}
			return nil
		}

		// 如果是文件，添加到 ZIP 并压缩内容
		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// 打开源文件
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 将源文件内容写入 ZIP 文件中
		_, err = io.Copy(zipFileWriter, srcFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to zip directory: %v", err)
	}

	return nil
}

// 打开zip包以供读取
func Zip_open(filename string) (*zip.ReadCloser, error) {
	return zip.OpenReader(filename)
}

// 把数据装入一个二进制字符串
func Pack(order binary.ByteOrder, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 从二进制字符串对数据进行解包
func Unpack(order binary.ByteOrder, data string) (interface{}, error) {
	var result []byte
	r := bytes.NewReader([]byte(data))
	err := binary.Read(r, order, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
