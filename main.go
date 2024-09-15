// Package main
// Wrote by yijian on 2024/09/15
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	dirs = flag.String("dirs", "", "Directory to handle.")

	createYearDir  = flag.Bool("create-year-dir", false, "Create directory by year.")
	createMonthDir = flag.Bool("create-month-dir", false, "Create directory by month.")
	createDayDir   = flag.Bool("create-day-dir", false, "Create directory by day.")
)

func main() {
	flag.Parse()
	if !checkParamers() {
		flag.Usage()
		os.Exit(1)
	}

	dirArray := strings.Split(*dirs, ",")
	for _, dir := range dirArray {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Read directory `%s` error: %s\n", path, err.Error())
				return err
			}
			if d.Type().IsRegular() {
				ext := filepath.Ext(path)
				if isImageOrViedoFormat(ext, "png", "jpeg", "jpg", "raw", "") {
					fmt.Println(path)
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

func isImageOrViedoFormat(filename string, formats ...string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, format := range formats {
		if ext == "."+format {
			return true
		}
	}
	return false
}
