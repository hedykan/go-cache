package cachemap

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheMap(t *testing.T) {
	b := NewCacheMap[int, int]()
	b.Set(0, 0, 5*time.Second)

	for i := 0; i < 10; i++ {
		b.Get(0)
		fmt.Println(b)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(6 * time.Second)
	fmt.Println(b)
}

func TestBuffSet(t *testing.T) {
	b := NewCacheMap[int, int]()

	for i := 0; i < 1000; i++ {
		go b.Set(1, i, time.Second)
	}
	time.Sleep(2 * time.Second)
}

func BenchmarkAdd(b *testing.B) {
	buff := NewCacheMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buff.Set(i, i, time.Second)
	}
}
