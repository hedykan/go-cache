package list

import (
	"container/list"
	"log"
	"sync"

	"github.com/hedykan/go-cache/types"
)

/*
顺序保存队列
- updateList: 保存要更新的副本，按更新顺序插入队列
- mu: 用来保证队列的插入是并发安全的
- retryCount: 重试次数，当保存失败时候，根据重试次数进行重试
- autoSave: 自动保存，如果自动保存开启，将会在每次插入时使用saveFrontNode
*/
type SaveList[T types.Saver] struct {
	updateList *list.List
	mu         *sync.Mutex
	retryCount int
	autoSave   bool
	saveStatus bool
}

func NewSaveList[T types.Saver]() *SaveList[T] {
	return &SaveList[T]{
		updateList: list.New(),
		mu:         &sync.Mutex{},
		retryCount: 0,
		autoSave:   true,
	}
}

// 设置重试次数
func (n *SaveList[T]) SetRetryCount(count int) {
	n.retryCount = count
}

// 设置自动保存
func (n *SaveList[T]) SetAutoSave(status bool) {
	n.autoSave = status
}

// 插入待更新内容
func (n *SaveList[T]) PushUpdate(val T) {
	n.mu.Lock()
	n.updateList.PushBack(val)
	n.mu.Unlock()

	// 如果开启了自动保存
	if n.autoSave {
		go n.SaveAllNode()
	}
}

// 手动保存所有数据
func (n *SaveList[T]) SaveAllNode() {
	// 已经在保存了
	if n.saveStatus {
		// log.Println("print save status", n.saveStatus)
		return
	}
	// 否则退出的时候进行更新
	n.saveStatus = true
	defer func() {
		n.saveStatus = false
		// log.Println("print save status", n.saveStatus)
	}()

	n.mu.Lock()
	defer n.mu.Unlock()

	retryCount := 0
	e := n.updateList.Front()
	for e != nil {
		if retryCount > n.retryCount {
			log.Println("retry count max, retry count:", retryCount, "n.retryCount:", n.retryCount)
			break
		}
		v := e.Value.(types.Saver)
		// log.Println("print save", v)
		err := v.Save()
		if err != nil {
			log.Println("save err:", err.Error(), ", retry count:", retryCount)
			retryCount += 1
			continue
		} else {
			retryCount = 0
		}

		// 继续下一条，直到处理完成
		n.updateList.Remove(e)
		e = n.updateList.Front()
	}
}
