package hashutil

import "testing"

func BenchmarkHashFromFileName(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testHash.HashFromFileName(testFilePath)
		if err != nil {
			b.Errorf("test HashFromFileName error =>%v", err)
		}
	}
}

func BenchmarkCheckpointsHashFromFileName_OnlyEntireFile(b *testing.B) {
	benchmarkCheckpointsHashFromFileName(b, 0, 0)
}

func BenchmarkCheckpointsHashFromFileName(b *testing.B) {
	benchmarkCheckpointsHashFromFileName(b, 1024, 10)
}

func BenchmarkCheckpointsHashFromFileName_WithDefaultChunkSize(b *testing.B) {
	benchmarkCheckpointsHashFromFileName(b, defaultChunkSize, 10)
}

func benchmarkCheckpointsHashFromFileName(b *testing.B, chunkSize int64, checkpointCount int) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testHash.CheckpointsHashFromFileName(testFilePath, chunkSize, checkpointCount)
		if err != nil {
			b.Errorf("benchmark test CheckpointsHashFromFileName error chunkSize=%d checkpointCount=%d =>%v", chunkSize, checkpointCount, err)
		}
	}
}
