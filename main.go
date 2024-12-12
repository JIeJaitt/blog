package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func printTree(root string, prefix string) {
	// 读取目录内容
	dirEntries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var entries []fs.DirEntry
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			// 如果你也想打印文件，则取消下一行的注释
			// fmt.Println(prefix + "├── " + entry.Name())
		} else {
			entries = append(entries, entry)
		}
	}

	// 遍历目录
	for i, entry := range entries {
		// 处理前缀，最后一个目录/文件与其他的区别对待
		entryPrefix := prefix + "├── "
		if i == len(entries)-1 {
			entryPrefix = prefix + "└── "
		}

		fmt.Println(entryPrefix + entry.Name())
		// 递归处理子目录，增加缩进
		nextPrefix := prefix + "│   "
		if i == len(entries)-1 {
			nextPrefix = prefix + "    "
		}

		printTree(filepath.Join(root, entry.Name()), nextPrefix)
	}
}

func main() {
	// 遍历/etc目录
	printTree("/etc", "")
}
