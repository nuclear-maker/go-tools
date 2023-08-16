package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	// 获取命令行参数
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("请提供目录地址作为参数")
		os.Exit(1)
	}

	// 递归处理目录下的文件
	for _, dirPath := range args {
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("访问文件或目录错误：%v\n", err)
				return nil
			}

			// 只处理文件，忽略目录
			if !info.IsDir() {
				// 检查文件编码并转换为 UTF-8
				convertedContent, err := convertToUtf8(path)
				if err != nil {
					fmt.Printf("转换文件失败：%v\n", err)
					return nil
				}

				// 写入转换后的内容
				err = os.WriteFile(path, convertedContent, info.Mode())
				if err != nil {
					fmt.Printf("写入文件失败：%v\n", err)
					return nil
				}

				fmt.Printf("已转换文件：%s\n", path)
			}

			return nil
		})

		if err != nil {
			fmt.Printf("遍历目录失败：%v\n", err)
		}
	}
}

// 检查文件编码并转换为 UTF-8
func convertToUtf8(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 检查文件编码
	if isGBK(content) {
		// 将 GBK 编码的内容转换为 UTF-8 编码
		reader := transform.NewReader(strings.NewReader(string(content)), simplifiedchinese.GBK.NewDecoder())
		utf8Content, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return utf8Content, nil
	}

	// 文件已经是 UTF-8 编码，无需转换
	return content, nil
}

// 检查文件内容是否为 GBK 编码
func isGBK(content []byte) bool {
	// 假设文件内容中包含非 ASCII 字符则判断为 GBK 编码
	for _, b := range content {
		if b > 127 {
			return true
		}
	}
	return false
}
