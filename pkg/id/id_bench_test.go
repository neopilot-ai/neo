package id_test

import (
	"testing"
)

func BenchmarkAscending(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Ascending()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDescending(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Descending()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConcurrentAscending(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Ascending()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkConcurrentDescending(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Descending()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkIDOrdering(b *testing.B) {
	ids := make([]string, 1000)
	var err error
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ids[i%1000], err = Ascending()
		if err != nil {
			b.Fatal(err)
		}
	}
}
