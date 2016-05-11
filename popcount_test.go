package hamt

import "testing"

func BenchmarkPopcount2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount_2(uint64(i))
	}
}

func BenchmarkPopcount3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount_3(uint64(i))
	}
}

func BenchmarkPopcnt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcnt(uint64(i))
	}
}
