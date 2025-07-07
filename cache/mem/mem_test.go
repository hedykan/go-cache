package mem

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedykan/go-cache/types"
)

// 每次更新可以指定更新方法
// 在保存函数中打印err，可以知道保存哪里出问题
type testType struct {
	Val        int
	updateFunc func() error
	types.BaseSave
}

func (t testType) Save() error {
	// time.Sleep(time.Duration(200-rand.Intn(50)) * time.Microsecond)
	if t.updateFunc != nil {
		err := t.updateFunc()
		return err
	}
	return nil
}

func TestCache(t *testing.T) {
	// node := NewCacheNode(testType{Val: -1, updateFunc: func() error {
	// 	log.Println("update")
	// 	err := errors.New("save error")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	return err
	// }})

	node := NewCacheNode(testType{Val: -1})
	// node.SetAutoSave(false)
	node.SetRetryCount(10)

	val := node.Get()
	fmt.Println(val)

	for i := 0; i < 1000; i++ {

		val.Val = i
		node.Update(val)
		val = node.Get()
		fmt.Println(val.Val)
	}

	node.SaveAllNode()

	time.Sleep(3 * time.Second)
}

func BenchmarkCache(b *testing.B) {
	node := NewCacheNode(testType{Val: 0})
	val := node.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.Val = i
		node.Update(val)
		val = node.Get()
	}
}
