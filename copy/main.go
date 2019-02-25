package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CopyDir copy srcPath dir to destPath dir but exclude srcPath dir,
// srcPath and destPath must exist.
func CopyDir(srcPath string, destPath string) error {

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
		if !f.IsDir() && !strings.HasSuffix(f.Name(), "_test.go") && !strings.Contains(f.Name(), "_autogen_") {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			copyFile(path, destNewPath)
		}
		return nil
	})

	return err
}

func copyFile(src, dest string) (w int64, err error) {

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
	err := CopyDir("/Users/test/srcPath", "/Users/test/destPath")
	if err != nil{
		panic(err)
	}
}

