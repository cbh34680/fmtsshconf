package main

import (
	"io"
	"os"
)

// PathExists ... https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go/12527546#12527546
func PathExists(name string) bool {

	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// CopyFile ... https://qiita.com/cotrpepe/items/93e4a072c249a933e795
func CopyFile(srcName, dstName string) error {

	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
