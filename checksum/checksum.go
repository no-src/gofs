package checksum

import (
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/gofs/util/jsonutil"
	"github.com/no-src/log"
)

// PrintChecksum calculate and print the checksum for file
func PrintChecksum(path string, chunkSize int64, checkpointCount int) error {
	hvs, err := hashutil.CheckpointsMD5FromFileName(path, chunkSize, checkpointCount)
	if err != nil {
		log.Error(err, "calculate file checksum error")
		return err
	}

	hvsJson, _ := jsonutil.MarshalIndent(hvs)
	log.Log(string(hvsJson))
	return err
}
