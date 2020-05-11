package zip

import (
	"os"
	"testing"
)

func TestZipUnzip(t *testing.T) {

	err := Zip("testdata/zipped.zip", "testdata/testDir")
	if err != nil {
		t.Errorf("Got error when zipping: %v", err)
	}

	err = Unzip("testdata/zipped.zip", "testdata/unzipped")
	if err != nil {
		t.Errorf("Got error when unzipping: %v", err)
	}

	_, err = os.Stat("testdata/unzipped/testDir/testFile1.txt")
	if err != nil {
		t.Errorf("Got error when stating testFile1: %v", err)
	}

	_, err = os.Stat("testdata/unzipped/testDir/innerDir/testFile2.txt")
	if err != nil {
		t.Errorf("Got error when stating testFile2: %v", err)
	}

	err = os.RemoveAll("testdata/unzipped")
	if err != nil {
		t.Errorf("Got error removing unzipped directory: %v", err)
	}

	err = os.Remove("testdata/zipped.zip")
	if err != nil {
		t.Errorf("Got error removing zip file: %v", err)
	}
}
