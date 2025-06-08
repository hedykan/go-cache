package cachemap

import (
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
	m        map[keyT]CacheMapNode[T]
	saveList *list.SaveList[T]
	mu       *sync.Mutex
}

/*
缓存表节点
用以实现保存相关元素，并且到期自动清除，访问自动续期
- val: 具体值的内容
- reset: 重置函数，每次访问之后会自动调用用以续期
*/
type CacheMapNode[T types.Saver] struct {
	val   T
	reset func()
}

// 新建缓存表
func NewCacheMap[keyT IBuffKey, T types.Saver]() *CacheMap[keyT, T] {
	saveList := list.NewSaveList[T]()
	bm := &CacheMap[keyT, T]{
		m:        make(map[keyT]CacheMapNode[T]),
		saveList: saveList,
		mu:       &sync.Mutex{},
	}
	return bm
}

// 从缓存表中读取元素
func (b CacheMap[keyT, T]) Get(key keyT) (T, bool) {
	val, ok := b.m[key]
	if ok {
		val.reset()
	}
	return val.val, ok
}

// 在缓存表中设置元素
func (b CacheMap[keyT, T]) Set(key keyT, val T, seg time.Duration) {
	node := CacheMapNode[T]{}
	t := time.NewTimer(seg)
	reset := func() {
		if !t.Stop() {
			select {
			case <-t.C:
			default:
			}
		}
		t.Reset(seg)
	}

	node.val = val
	node.reset = reset

	b.mu.Lock()
	b.m[key] = node
	b.mu.Unlock()

	// 插入数据库更新
	b.saveList.PushUpdate(node.val)

	go func() {
		<-t.C
		b.Delete(key)
	}()
}

// 删除指定元素
func (b CacheMap[keyT, T]) Delete(key keyT) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.m, key)
}
