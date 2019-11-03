package mytimewheel

import (
	"container/list"
	"sync/atomic"
	"unsafe"
)

type Timer struct {
	expire int64
	task   func()
	e      *list.Element
	bucket unsafe.Pointer
}

func newTimer(expire int64, task func()) *Timer {
	return &Timer{
		expire: expire,
		task:   task,
	}
}

func (t *Timer) Stop() {
	var suc bool
	for b := t.getBucket(); b != nil && !suc; b = t.getBucket() {
		suc = b.delTimer(t)
	}
}
func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.bucket))
}

func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.bucket, unsafe.Pointer(b))
}
