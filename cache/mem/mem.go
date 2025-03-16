package mem

import (
	slist "github.com/hedykan/go-cache/save/list"
	"github.com/hedykan/go-cache/types"
)

// 缓存接口
type ICacheValue interface {
	types.Saver
}

/*
缓存节点
缓存节点包含要保存的内容，一个更新队列，以及一个队列锁
- val: 作为最新的内容的值提供查询，这个内容需要实现一个保存接口，用来执行保存操作
- updateList: 顺序保存队列，将每次更新的值插入进顺序保存队列，将会自动进行保存
*/
type CacheNode[T ICacheValue] struct {
	val        T
	updateList *slist.SaveList[T]
}

// 新建缓存节点
func NewCacheNode[T ICacheValue](val T) *CacheNode[T] {
	return &CacheNode[T]{
		val:        val,
		updateList: slist.NewSaveList[T](),
	}
}

// 设置重试次数
func (n *CacheNode[T]) SetRetryCount(count int) {
	n.updateList.SetRetryCount(count)
}

// 设置自动保存
func (n *CacheNode[T]) SetAutoSave(status bool) {
	n.updateList.SetAutoSave(status)
}

// 获取节点内容
func (n *CacheNode[T]) Get() T {
	return n.val
}

// 更新节点内容
func (n *CacheNode[T]) Update(val T) {
	n.val = val
	n.updateList.PushUpdate(val)
}

func (n *CacheNode[T]) SaveAllNode() {
	n.updateList.SaveAllNode()
}
