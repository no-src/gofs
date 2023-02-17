package hashutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	compareTestFile = "./compare_test.go"
	compareGoFile   = "./compare.go"
)

func TestCompare(t *testing.T) {
	tempFile, err := os.CreateTemp("", filepath.Clean(compareTestFile))
	if err != nil {
		t.Errorf("create temp file error => %v", err)
		return
	}
	defer tempFile.Close()

	testFileData, err := os.ReadFile(compareTestFile)
	if err != nil {
		t.Errorf("read file error => %v", err)
		return
	}
	_, err = tempFile.Write(testFileData[:len(testFileData)/2])
	if err != nil {
		t.Errorf("write file error => %v", err)
		return
	}
	testCases := []struct {
		name     string
		source   string
		dest     string
		destSize int64
		equal    bool
	}{
		{"zero dest size", compareTestFile, compareTestFile, 0, false},
		{"the same files", compareTestFile, compareTestFile, 1, true},
		{"the different files", compareTestFile, compareGoFile, 1, false},
		{"only part of it is the same", compareTestFile, tempFile.Name(), 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src, err := os.Open(tc.source)
			if err != nil {
				t.Errorf("open file error => %v", err)
				return
			}
			defer src.Close()
			srcStat, err := src.Stat()
			if err != nil {
				t.Errorf("get file stat error => %v", err)
				return
			}
			var offset int64
			actual := Compare(100, 10, src, srcStat.Size(), tc.dest, tc.destSize, &offset)
			if actual != tc.equal {
				t.Errorf("expect get %v but actual get %v", tc.equal, actual)
			}
		})
	}
}

func TestQuickCompare(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name          string
		forceChecksum bool
		sourceSize    int64
		destSize      int64
		sourceModTime time.Time
		destModTime   time.Time
		equal         bool
	}{
		{"force checksum", true, 10, 10, now, now, false},
		{"force checksum with different size", true, 1, 2, time.Now(), time.Now(), false},
		{"equal", false, 10, 10, now, now, true},
		{"size not equal", false, 1, 2, now, now, false},
		{"time not equal", false, 10, 10, now, now.Add(time.Second), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := QuickCompare(tc.forceChecksum, tc.sourceSize, tc.destSize, tc.sourceModTime, tc.destModTime)
			if actual != tc.equal {
				t.Errorf("expect get %v but actual get %v", tc.equal, actual)
			}
		})
	}
}

func TestCompareHashValues(t *testing.T) {
	f, err := os.Open(compareTestFile)
	if err != nil {
		t.Errorf("open file error => %v", err)
		return
	}
	defer f.Close()

	currentTestFileStat, err := f.Stat()
	if err != nil {
		t.Errorf("get file stat error => %v", err)
		return
	}

	currentTestFileHash, err := HashFromFile(f)
	if err != nil {
		t.Errorf("get file hash error => %v", err)
		return
	}

	testCases := []struct {
		name       string
		dstPath    string
		sourceSize int64
		sourceHash string
		chunkSize  int64
		hvs        HashValues
		equal      bool
	}{
		{"zero source size", compareTestFile, 0, currentTestFileHash, 100, nil, false},
		{"the same files without hvs", compareTestFile, currentTestFileStat.Size(), currentTestFileHash, 100, nil, true},
		{"the different files without hvs", compareGoFile, currentTestFileStat.Size(), currentTestFileHash, 100, nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := CompareHashValues(tc.dstPath, tc.sourceSize, tc.sourceHash, tc.chunkSize, tc.hvs)
			if actual != tc.equal {
				t.Errorf("expect get %v but actual get %v", tc.equal, actual)
			}

		})
	}
}
