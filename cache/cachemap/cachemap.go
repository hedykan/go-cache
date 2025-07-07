package cachemap

import (
	"log"
	"sync"
	"time"

	"github.com/hedykan/go-cache/save/list"
	"github.com/hedykan/go-cache/types"
)

type IBuffKey interface {
	int | int8 | int16 | int32 | int64 | string
}

/*
缓存表
实现了一个到期定时清除的map
- m: 缓存表，用以存储值
- mu: 表插入锁，用以保证表插入的线程安全性
*/
type CacheMap[keyT IBuffKey, T types.Saver] struct {
	m         map[keyT]*CacheMapNode[T]
	saveList  *list.SaveList[T]
	mu        *sync.Mutex
	resetTime time.Duration // 重置时间
}

/*
缓存表节点
用以实现保存相关元素，并且到期自动清除，访问自动续期
- val: 具体值的内容
- reset: 重置函数，每次访问之后会自动调用用以续期
*/
type CacheMapNode[T types.Saver] struct {
	val        T
	updateTime time.Time // 更新时间
}

// 更新时间
func (node *CacheMapNode[T]) Reset() {
	node.updateTime = time.Now()
	log.Println("update cache", node.updateTime)
}

// 获取更新时间
func (node *CacheMapNode[T]) GetUpdateTime() time.Time {
	return node.updateTime
}

// 新建缓存表
func NewCacheMap[keyT IBuffKey, T types.Saver](cacheTime time.Duration) *CacheMap[keyT, T] {
	saveList := list.NewSaveList[T]()
	bm := &CacheMap[keyT, T]{
		m:         make(map[keyT]*CacheMapNode[T]),
		saveList:  saveList,
		mu:        &sync.Mutex{},
		resetTime: cacheTime,
	}
	go bm.reset()
	return bm
}

func (b CacheMap[keyT, T]) reset() {
	tick := time.NewTicker(b.resetTime)
	for {
		<-tick.C
		for key := range b.m {

			b.mu.Lock()
			val := b.m[key]
			b.mu.Unlock()

			if time.Since(val.GetUpdateTime()) > b.resetTime {
				log.Println("start clear cache", key)
				b.Delete(key)
			}
		}
	}
}

// 从缓存表中读取元素，读时更新
func (b CacheMap[keyT, T]) Get(key keyT) (T, bool) {
	b.mu.Lock()
	val, ok := b.m[key]
	b.mu.Unlock()
	if ok {
		updateTime := val.GetUpdateTime()
		// 如果超时，清除内容
		if time.Since(updateTime) > b.resetTime {
			log.Println("clear cache", time.Since(updateTime), updateTime)
			b.Delete(key)
			return val.val, false
		}

		val.Reset()
	}
	return val.val, ok
}

// 在缓存表中设置元素
func (b CacheMap[keyT, T]) Set(key keyT, val T, seg time.Duration) {
	node := CacheMapNode[T]{}
	node.val = val
	node.updateTime = time.Now()

	b.mu.Lock()
	b.m[key] = &node
	b.mu.Unlock()

	// 插入数据库更新
	b.saveList.PushUpdate(node.val)
}

// 删除指定元素
func (b CacheMap[keyT, T]) Delete(key keyT) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.m, key)
}
