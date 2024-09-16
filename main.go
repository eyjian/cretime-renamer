// Package main
// Wrote by yijian on 2024/09/15
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	// 指定文件目录，多个目录间已逗号分割
	dirs = flag.String("dirs", "", "Directory to handle.")

	createYearDir  = flag.Bool("create-year-dir", false, "Create directory by year.")
	createMonthDir = flag.Bool("create-month-dir", false, "Create directory by month.")
	createDayDir   = flag.Bool("create-day-dir", false, "Create directory by day.")

	// 指定需要处理的文件名后缀，如果为空表示处理所有的文件
	suffixes = flag.String("suffixes", "", "File name suffixes that needs to be processed.")
)

func main() {
	flag.Parse()
	if !checkParamers() {
		flag.Usage()
		os.Exit(1)
	}

	// 得到目录数组
	dirArray := strings.Split(*dirs, ",")

	// 处理所有目录
	for _, dir := range dirArray {
		// 遍历指定目录
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Read directory `%s` error: %s\n", path, err.Error())
				return err
			}
			if d.Type().IsRegular() {
				ext := filepath.Ext(path)
				if needProcess(ext) {
					fi, err := d.Info()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Stat file `%s` error: %s", path, err.Error())
					} else {
						rename(path, ext, fi)
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Walk directory `%s` error: %s\n", dir, err.Error())
		}
	}
}

func checkParamers() bool {
	if *dirs == "" {
		fmt.Fprintf(os.Stderr, "Parameter -dirs is not set.\n")
		return false
	}

	return true
}

func needProcess(filename string) bool {
	if *suffixes == "" {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filename))
	suffixArray := strings.Split(*suffixes, ",")
	for _, suffix := range suffixArray {
		if ext == "."+suffix {
			return true
		}
	}
	return false
}

func rename(path, ext string, fi fs.FileInfo) {
	dir := filepath.Dir(path)
	newPath, err := getNewFilepath(fi, ext, dir, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Rename file `%s` error: %s\n", path, err.Error())
		return
	}

	idx := 1
	for {
		exists, err := PathExists(newPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Rename file `%s` to `%s` error: %s\n", path, newPath, err.Error())
			return
		}
		if !exists {
			break
		} else {
			newPath, err = getNewFilepath(fi, ext, dir, idx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Rename file `%s` error: %s\n", path, err.Error())
				return
			}
			idx++
		}
	}

	err = os.Rename(path, newPath)
	if err == nil {
		fmt.Fprintf(os.Stdout, "Rename file `%s` to `%s` ok\n", path, newPath)
	} else {
		fmt.Fprintf(os.Stderr, "Rename file `%s` to `%s` error: %s\n", path, newPath, err.Error())
	}
}

func getNewFilepath(fi fs.FileInfo, ext, dir string, idx int) (string, error) {
	fileDir := dir

	if *createYearDir {
		year := fi.ModTime().Format("20060102")
		fileDir := fmt.Sprintf("%s%c%s", dir, filepath.Separator, year)
		exists, err := DirExists(fileDir)
		if err != nil {
			return "", err
		}
		if !exists {
			err = os.Mkdir(fileDir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Create directory `%s` error: %s\n", fileDir, err.Error())
				return "", err
			}
		}
	}

	if idx < 1 {
		return fmt.Sprintf("%s%c%s%s", fileDir, filepath.Separator, fi.ModTime().Format("20060102150405"), ext), nil
	} else {
		return fmt.Sprintf("%s%c%s-%02d%s", fileDir, filepath.Separator, fi.ModTime().Format("20060102150405"), idx, ext), nil
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// 其他错误
	return false, err
}

func DirExists(path string) (bool, error) {
	st, err := os.Stat(path)
	if err == nil {
		if st.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("not a directory")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// 其他错误
	return false, err
}
