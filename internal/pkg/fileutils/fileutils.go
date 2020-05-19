package fileutils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rajch/kutti/internal/pkg/kuttilog"
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
func CopyFile(sourcepath string, destpath string, buffersize int64, overwrite bool) error {

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

	if !overwrite {
		_, err = os.Stat(destpath)
		if err == nil {
			return fmt.Errorf("file %s already exists", destpath)
		}
	}

	destination, err := os.Create(destpath)
	if err != nil {
		return err
	}
	defer destination.Close()

	// Debug log output
	if kuttilog.V(4) {
		kuttilog.Printf(4, "Copying %s to %s:\n", sourcepath, destpath)
	}

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

		// if kuttilog.V(4) {
		// 	fmt.Print(".")
		// }
	}

	// if kuttilog.V(4) {
	// 	fmt.Println(".")
	// }

	return err
}

// DownloadFile downloads a file from a url.
func DownloadFile(url string, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	tmpfilepath := filepath + ".download"
	out, err := os.Create(tmpfilepath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		out.Close()
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	if err = os.Rename(tmpfilepath, filepath); err != nil {
		return err
	}
	return nil
}

// RemoveFile deletes a file.
func RemoveFile(filepath string) error {
	return os.Remove(filepath)
}
