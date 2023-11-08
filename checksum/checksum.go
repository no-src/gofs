package checksum

import (
	"github.com/no-src/gofs/logger"
	"github.com/no-src/nsgo/hashutil"
	"github.com/no-src/nsgo/jsonutil"
)

// PrintChecksum calculate and print the checksum for file
func PrintChecksum(path string, chunkSize int64, checkpointCount int, algorithm string, logger *logger.Logger) error {
	hash, err := hashutil.NewHash(algorithm)
	if err != nil {
		logger.Error(err, "init hash component error")
		return err
	}
	hvs, err := hash.CheckpointsHashFromFileName(path, chunkSize, checkpointCount)
	if err != nil {
		logger.Error(err, "calculate file checksum error")
		return err
	}

	hvsJson, _ := jsonutil.MarshalIndent(hvs)
	logger.Log(string(hvsJson))
	return err
}
