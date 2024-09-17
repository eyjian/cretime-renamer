// Package main
// Wrote by yijian on 2024/09/15
// 本工具用来处理手机和相机导出的照片和视频，使用相片和视频的创建时间作为文件名。
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// 指定文件目录，多个目录间已逗号分割
	dirs = flag.String("dirs", "", "Directory to handle.")

	// 是否按年饭创建目录
	createYearDir = flag.Bool("create-year-dir", false, "Create directory by year.")

	// 是否按月份创建目录
	// 如果 create-month-dir 为 true，则强制 create-year-dir 也为 true
	createMonthDir = flag.Bool("create-month-dir", false, "Create directory by month.")

	// 指定需要处理的文件名后缀，如果为空表示处理所有的文件
	suffixes = flag.String("suffixes", "", "File name suffixes that needs to be processed.")

	// 是否在同级目录下创建年份目录
	// 仅当 create-year-dir 或 create-month-dir 为 true 时有作用
	// 如果为 false，则年份目录会创建在同文件同级的目录
	siblingDirectory = flag.Bool("sibling-dir", false, "Create directory by year in the sibling directory.")

	// 需要忽略的目录列表
	ignoreDirs = flag.String("ignore-dirs", "", "Directory ignored.")

	// 跳过日期目录（仅 create-year-dir 或 create-month-dir 均为 true 时有效）
	skipDateDir = flag.Bool("skip-date-dir", true, "Skip date directory.")
)

func main() {
	flag.Parse()
	if !checkParamers() {
		flag.Usage()
		os.Exit(1)
	}

	// 得到目录数组
	dirArray := strings.Split(*dirs, ",")

	// 被忽略的目录数组
	ignoreDirArray := strings.Split(*ignoreDirs, ",")

	// 处理所有目录
	for _, dir := range dirArray {
		// 遍历指定目录
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Read directory `%s` error: %s\n", path, err.Error())
				return err
			}
			if d.Type().IsDir() {
				tmpDir := filepath.Join(path, "xyz-###-123")
				leafDir := filepath.Base(filepath.Dir(tmpDir))
				if isSkippedDir(leafDir) {
					fmt.Fprintf(os.Stderr, "Directory `%s` is skipped\n", path)
					return filepath.SkipDir
				}
				if isIgnoredDirs(ignoreDirArray, path) {
					fmt.Fprintf(os.Stderr, "Directory `%s` is ignored\n", path)
					return filepath.SkipDir
				}
			} else if d.Type().IsRegular() {
				if notNeedProcess(path) {
					fmt.Fprintf(os.Stderr, "Path `%s` need not been processed\n", path)
				} else {
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

func notNeedProcess(path string) bool {
	shortname := GetFileNameWithoutExtension(path)
	return IsValidYYYYMMDDhhmmss(shortname)
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
		err = os.Chtimes(newPath, fi.ModTime(), fi.ModTime())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Chtimes file `%s` error: %s\n", newPath, err.Error())
		}
		err = os.Chmod(newPath, fi.Mode())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Chmod file `%s` error: %s\n", newPath, err.Error())
		}
	} else {
		fmt.Fprintf(os.Stderr, "Rename file `%s` to `%s` error: %s\n", path, newPath, err.Error())
	}
}

func getNewFilepath(fi fs.FileInfo, ext, dir string, idx int) (string, error) {
	fileDir := dir

	if *createYearDir {
		year := fi.ModTime().Format("2006")
		baseDir := dir

		if *siblingDirectory {
			baseDir = filepath.Dir(baseDir)
		}
		fileDir = fmt.Sprintf("%s%c%s", baseDir, filepath.Separator, year)
		exists, err := DirExists(fileDir)
		if err != nil {
			return "", err
		}
		if !exists {
			err = os.MkdirAll(fileDir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Create directory `%s` error: %s\n", fileDir, err.Error())
				return "", err
			}
		}

		if *createMonthDir {
			yearMonth := fi.ModTime().Format("200601")
			fileDir = fmt.Sprintf("%s%c%s%c%s", baseDir, filepath.Separator, year, filepath.Separator, yearMonth)
			exists, err = DirExists(fileDir)
			if err != nil {
				return "", err
			}
			if !exists {
				err = os.MkdirAll(fileDir, 0755)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Create directory `%s` error: %s\n", fileDir, err.Error())
					return "", err
				}
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

func IsValidYYYYMMDDhhmmss(s string) bool {
	if len(s) != 14 {
		return false
	}

	t, err := time.Parse("20060102150405", s)
	return err == nil && t.Year() >= 1978
}

func IsValidYYYYMMDD(s string) bool {
	var err error
	var t time.Time

	l := len(s)
	if l == 4 {
		t, err = time.Parse("2006", s)
	} else if l == 6 {
		t, err = time.Parse("200601", s)
	} else {
		t, err = time.Parse("20060102", s)
	}

	return err == nil && t.Year() >= 1978
}

func isIgnoredDirs(ignoredDirArray []string, dir string) bool {
	for _, ignoredDir := range ignoredDirArray {
		if ignoredDir == dir {
			return true
		}
	}

	return false
}

func isSkippedDir(dir string) bool {
	if IsValidYYYYMMDD(dir) {
		if *createMonthDir || *createYearDir {
			return *skipDateDir
		}
	}
	return false
}

func GetFileNameWithoutExtension(path string) string {
	base := filepath.Base(path)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	return fileName
}
