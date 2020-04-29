package fileutils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// ChecksumFile calculates an SHA256 checksum of a file
func ChecksumFile(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	result := fmt.Sprintf("%x", h.Sum(nil))
	return result, nil
}

// CopyFile copies a file
func CopyFile(sourcepath string, destpath string, buffersize int64) error {

	sourceFileStat, err := os.Stat(sourcepath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourcepath)
	}

	source, err := os.Open(sourcepath)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(destpath)
	if err == nil {
		return fmt.Errorf("file %s already exists", destpath)
	}

	destination, err := os.Create(destpath)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, buffersize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
