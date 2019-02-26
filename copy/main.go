package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CopyDir copy srcPath dir to destPath dir but exclude srcPath dir,
// srcPath and destPath must exist.
// copy file that name match matchPattern and not match excludePattern
func CopyDir(srcPath, destPath, matchPattern, excludePattern string) error {

	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath is not a dir")
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo is not a dir")
			return e
		}
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		match, errmc := regexp.MatchString(matchPattern, f.Name())
		if errmc != nil {
			return errors.New(errmc.Error())
		}

		exclude, errex := regexp.MatchString(excludePattern, f.Name())
		if errex != nil {
			return errors.New(errmc.Error())
		}

		if !f.IsDir() && match && !exclude {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			_, err := CopyFile(path, destNewPath)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// CopyFile copy src file to dest
func CopyFile(src, dest string) (w int64, err error) {

	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	destSplitPathDirs := strings.Split(dest, "/")

	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, e := pathExists(destSplitPath)
			if e != nil {
				err = errors.New(e.Error())
				return
			}
			if b == false {
				e := os.Mkdir(destSplitPath, os.ModePerm)
				if e != nil {
					err = errors.New(e.Error())
					return
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func pathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	err := CopyDir("/Users/test/copy/a", "/Users/test/copy/b", "(.go)$", "(_notCopy_)")
	if err != nil {
		fmt.Println(err)
	}

	w, e := CopyFile("/Users/test/copy/a/told.go", "/Users/test/copy/b/new.go")
	fmt.Println(w)
	fmt.Println(e)
}
