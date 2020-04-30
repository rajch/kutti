package fileutils

import "testing"

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

	err = CopyFile("fileutils_test.go", "deletethis_test.xxx", 1000)
	if err != nil {
		t.Logf("Copyfile failed with error:%v", err)
		t.FailNow()
	}

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
