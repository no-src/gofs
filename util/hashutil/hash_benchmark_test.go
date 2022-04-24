package hashutil

import "testing"

func BenchmarkMD5FromFileName(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MD5FromFileName(testFilePath)
		if err != nil {
			b.Errorf("test MD5FromFileName error =>%v", err)
		}
	}
}

func BenchmarkCheckpointsMD5FromFileName_OnlyEntireFile(b *testing.B) {
	benchmarkCheckpointsMD5FromFileName(b, 0, 0)
}

func BenchmarkCheckpointsMD5FromFileName(b *testing.B) {
	benchmarkCheckpointsMD5FromFileName(b, 1024, 10)
}

func BenchmarkCheckpointsMD5FromFileName_WithDefaultChunkSize(b *testing.B) {
	benchmarkCheckpointsMD5FromFileName(b, defaultChunkSize, 10)
}

func benchmarkCheckpointsMD5FromFileName(b *testing.B, chunkSize int64, checkpointCount int) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CheckpointsMD5FromFileName(testFilePath, chunkSize, checkpointCount)
		if err != nil {
			b.Errorf("benchmark test CheckpointsMD5FromFileName error chunkSize=%d checkpointCount=%d =>%v", chunkSize, checkpointCount, err)
		}
	}
}
