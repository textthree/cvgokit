package filekit

import (
	"bufio"
	"cvgo/kit/strkit"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 获取当前路径（运行命令时所在的目录，切换到哪个目录去运行当前就是哪个目录）
func Getwd() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return currentDir
}

// 获取一个绝对路径所属目录
func Dir(absolutePath string) string {
	return filepath.Dir(absolutePath)
}

// 判断文件/文件夹是否存在
func PathExists(absolutePath string) (bool, error) {
	_, err := os.Stat(absolutePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断是否是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断是否是文件而不是目录
func IsFile(path string) bool {
	return !IsDir(path)
}

// 根据文件名获取后缀
func GetSuffix(filename string) string {
	s := strkit.Explode(".", filename)
	last := s[len(s)-1]
	if last == filename {
		return ""
	}
	return last
}

// 扫描指定目录下所有文件及文件夹
// dir 指定要扫描的目录
// return：
// files 文件数组
// dirs  文件夹数组
func Scandir(dir string) (files []string, dirs []string) {
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			// 扫描到的 ./ (自身) 不要
			if path != dir {
				dirs = append(dirs, path)
			}
			return nil
		}
		files = append(files, path)
		return nil
	})
	return
}

// 创建文件，覆盖创建
func createFile(filepath string, content string) {
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(content))
	}
}

// 读取文件内容
func readFile(filepath string) string {
	ret := ""
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		contentByte, _ := io.ReadAll(f)
		ret = string(contentByte)
	}
	return ret
}

// 上传文件保存到本地
func UploadFile(file *multipart.FileHeader) {
	// 输出文件信息
	fmt.Printf("Uploaded File: %+v\n", file.Filename)
	fmt.Printf("File Size: %+v\n", file.Size)
	fmt.Printf("MIME Header: %+v\n", file.Header)
}

// 判断是否图片类型
func IsImage(file *multipart.FileHeader) bool {
	// 获取文件的 MIME 类型
	mimeType := file.Header.Get("Content-Type")
	// 检查 MIME 类型是否是允许的图片类型
	return strings.HasPrefix(mimeType, "image/")
}

// 删除给定路径的文件。
// 此函数可以删除控目录，但不能删除非空目录
// filePath 文件的绝对路径或相对路径
func DeleteFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// CreatePath 创建文件或目录
// 如果路径以 "/" 结尾，将被认为是目录；否则认为是文件。
// overwite 如果目标已经存在，是否覆盖
func CreatePath(path string, overwite ...bool) error {
	// 判断路径是否是目录
	isDir := filepath.Ext(path) == ""
	force := false
	if len(overwite) > 0 {
		force = overwite[0]
	}
	if !force {
		exists, err := PathExists(path)
		if err != nil || exists {
			return errors.New("创建失败，路径已存在：" + path)
		}
	}
	if isDir {
		// 创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	} else {
		// 创建文件所在的目录
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating parent directory: %v", err)
		}

		// 创建文件
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()
	}

	return nil
}

// CopyFile 将 src 文件复制到 dst，如果目录不在会自动创建
// 如果目标文件已存在默认会跳过，传递 override 后会覆盖
func CopyFile(src, dst string, override ...bool) error {
	force := false
	if len(override) > 0 {
		force = override[0]
	}
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		fmt.Println("CopyFile 打开源文件失败", err)
		return err
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	destDir := filepath.Dir(dst)
	if err := EnsureDirExists(destDir); err != nil {
		fmt.Println("EnsureDirExists Error:", err)
		return err
	}

	// 检查目标文件是否已存
	ext, err := PathExists(dst)
	if ext && !force {
		return errors.New("目标文件已经存在，跳过创建 " + dst)
	}

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("CopyFile 创建目标文件失败", err)
		return err
	}
	defer destFile.Close()

	// 将源文件内容复制到目标文件
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		fmt.Println("CopyFile 将源文件内容复制到目标文件失败", err)
		return err
	}

	// 确保写入完成
	err = destFile.Sync()
	if err != nil {
		return err
	}

	// 复制文件权限
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

// CopyFiles 将 srcDir 目录下的所有文件复制到 destDir 目录下
func CopyFiles(srcDir, dstDir string, override ...bool) error {
	// 检查目标目录是否存在，不存在则创建
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to add destination directory: %v", err)
		}
	}

	// 遍历源目录中的所有文件
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("遍历出错:", err)
			return err
		}

		// 忽略目录，只复制文件
		if !info.IsDir() {
			// 生成目标文件的路径
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				fmt.Printf("Failed to get relative path: %+v\n", err)
				return err
			}
			destPath := filepath.Join(dstDir, relPath)

			// 创建目标文件所在的目录结构
			destDirPath := filepath.Dir(destPath)
			if _, err := os.Stat(destDirPath); os.IsNotExist(err) {
				err = os.MkdirAll(destDirPath, os.ModePerm)
				if err != nil {
					fmt.Printf("Failed to add destination directory: %+v\n", err)
					return err
				}
			}

			// 复制文件
			err = CopyFile(path, destPath, override...)
			if err != nil {
				fmt.Printf("Failed to copy file: %+v\n", err)
				//	return err
			}
		}
		return nil
	})

	return err
}

// 移动目录，目标目录不存在会创建
func MoveDir(src, dst string) error {
	// 检查源路径是否存在
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", src)
	}
	// 确保目标目录存在，不存在就创建
	EnsureDirExists(dst)
	dst = filepath.Join(dst, filepath.Base(src))
	// 移动
	err = os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to move %s to %s: %v", src, dst, err)
	}
	return nil
}

// MoveFiles 将 srcDir 目录下的所有文件移动到 destDir 目录下
func MoveFiles(srcDir, destDir string) error {
	// 检查 srcDir 是否存在
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", srcDir)
	}

	// 检查 destDir 是否存在，如果不存在则创建它
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating destination directory: %v", err)
		}
	}

	// 读取 srcDir 目录下的所有文件
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("error reading source directory: %v", err)
	}

	// 遍历文件列表并移动它们到目标目录
	for _, file := range files {
		// 构建源文件路径和目标文件路径
		srcFilePath := filepath.Join(srcDir, file.Name())
		destFilePath := filepath.Join(destDir, file.Name())

		// 移动文件
		err := os.Rename(srcFilePath, destFilePath)
		if err != nil {
			return fmt.Errorf("error moving file %s: %v", file.Name(), err)
		}
	}
	return nil
}

// Rename 重命名目录或文件
func Rename(oldDir, newDir string) error {
	// 使用 os.Rename 进行重命名操作
	err := os.Rename(oldDir, newDir)
	if err != nil {
		return fmt.Errorf("failed to rename directory from %s to %s: %v", oldDir, newDir, err)
	}
	return nil
}

// 判断文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 创建文件并写入内容
func FilePutContents(filePath, content string) {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	EnsureDirExists(dir)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("创建文件 "+filePath+" 失败", err)
		return
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("往文件 "+filePath+" 写入内容失败", err)
		return
	}
}

// 读取文件的内容
func FileGetContents(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// AppendToFile 读取文件内容，追加新内容并写回文件
func FileAppendContent(filepath string, contentToAppend string) error {
	// 读取文件内容
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("无法读取文件: %w", err)
	}
	// 追加新内容
	newContent := string(fileContent) + contentToAppend
	// 将新内容写回文件
	err = os.WriteFile(filepath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("无法写入文件: %w", err)
	}
	return nil
}

// AddContentAboveLine 匹配指定行（会对行去除空白字符再匹配），在它前面添加内容。只会匹配一次
func AddContentAboveLine(filepath, matchLine, content string) error {
	alreadyMatched := false
	// 打开文件进行读取
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 创建一个新的内容容器
	var newContent strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 匹配指定行
		if !alreadyMatched && (strings.HasPrefix(strkit.RemoveSpace(line), strkit.RemoveSpace(matchLine))) {
			newContent.WriteString(content)
		}

		// 将当前行添加到新的内容中
		newContent.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %w", err)
	}

	// 将新的内容写回文件
	err = os.WriteFile(filepath, []byte(newContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("无法写入文件: %w", err)
	}

	return nil
}

// AddContentAboveLine 匹配指定行内容（会对行去除空白字符再匹配），在它后面添加内容。只会匹配一次
func AddContentUnderLine(filepath, matchLine, content string) error {
	alreadyMatched := false
	// 打开文件进行读取
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 创建一个新的内容容器
	var newContent strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 先将当前行添加到新的内容中
		newContent.WriteString(line + "\n")

		// 然后如果匹配到指定行，再插入一行内容
		if !alreadyMatched && (strings.HasPrefix(strkit.RemoveSpace(line), strkit.RemoveSpace(matchLine))) {
			newContent.WriteString(content + "\n")
			alreadyMatched = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %w", err)
	}

	// 将新的内容写回文件
	err = os.WriteFile(filepath, []byte(newContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("无法写入文件: %w", err)
	}

	return nil
}
