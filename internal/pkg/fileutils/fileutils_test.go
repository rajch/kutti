package fileutils

import (
	"os"
	"testing"
)

func TestChecksum(t *testing.T) {
	result, err := ChecksumFile("fileutils_test.go")
	if err != nil {
		t.Logf("Checksum failed with error:%v", err)
		t.FailNow()
	}
	t.Logf("Checksum is '%v'", result)
}

func TestCopyFile(t *testing.T) {
	sourceresult, err := ChecksumFile("fileutils_test.go")
	if err != nil {
		t.Logf("Checksum source failed with error:%v", err)
		t.FailNow()
	}
	t.Logf("Source Checksum is '%v'", sourceresult)

	err = CopyFile("fileutils_test.go", "deletethis_test.xxx", 1000, true)
	if err != nil {
		t.Logf("Copyfile failed with error:%v", err)
		t.FailNow()
	}

	defer os.Remove("deletethis_test.xxx")

	destresult, err := ChecksumFile("deletethis_test.xxx")
	if err != nil {
		t.Logf("Checksum of destination failed with error:%v", err)
		t.FailNow()
	}
	t.Logf("Destination Checksum is '%v'", destresult)

	if destresult != sourceresult {
		t.Log("Source and destination checksums don't match. Copy was faulty.")
		t.FailNow()
	}

}

func TestDownloadFile(t *testing.T) {
	const dfilename = "soedit-1.0.tar.gz"
	err := DownloadFile("https://github.com/rajch/soedit/releases/download/v1.0/soedit-1.0.tar.gz", dfilename)
	if err != nil {
		t.Logf("Downloadfile failed with error:%v\n", err)
		t.FailNow()
	}

	defer os.Remove(dfilename)

	destresult, err := ChecksumFile(dfilename)
	if err != nil {
		t.Logf("Checksum of destination failed with error:%v", err)
		t.FailNow()
	}
	t.Logf("Destination Checksum is '%v'", destresult)

	if destresult != "1f7960f2a6629b7af53c2cd1e1a505691573f2ba52641f23ca8cbf4814aa3526" {
		t.Log("Checksum don't match. Download was faulty.")
		t.FailNow()
	}

}
