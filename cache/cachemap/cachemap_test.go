package cachemap

import (
	"log"

	"testing"
	"time"
)

type SaveInt int

func (i SaveInt) Save() error {
	log.Println("save", i)
	return nil
}

func TestCacheMap(t *testing.T) {
	b := NewCacheMap[int, SaveInt](2 * time.Second)
	b.Set(0, 0, 5*time.Second)

	for i := 0; i < 10; i++ {
		b.Get(0)
		log.Println(b)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(6 * time.Second)
	log.Println(b)
}

func TestBuffSet(t *testing.T) {
	b := NewCacheMap[int, SaveInt](10 * time.Second)

	for i := 0; i < 1000; i++ {
		go b.Set(1, SaveInt(i), time.Second)
	}
	time.Sleep(2 * time.Second)
}

func BenchmarkAdd(b *testing.B) {
	buff := NewCacheMap[int, SaveInt](10 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buff.Set(i, SaveInt(i), time.Second)
	}
}
