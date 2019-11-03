package mytimewheel

import (
	"container/list"
	"sync"
)

type bucket struct {
	sync.Mutex
	// timers []*timer
	timers *list.List
}

func newBucket() *bucket {
	return &bucket{
		timers: list.New(),
	}
}

func (b *bucket) delTimer(t *Timer) bool {
	b.Lock()
	defer b.Unlock()
	if t.getBucket() != b {
		return false
	}
	b.timers.Remove(t.e)
	t.setBucket(nil)
	t.e = nil
	return true
}

func (b *bucket) addTimer(t *Timer) {
	b.Lock()
	defer b.Unlock()
	e := b.timers.PushBack(t)
	t.setBucket(b)
	t.e = e
}

func (b *bucket) getTimers() []*Timer {
	b.Lock()
	defer b.Unlock()
	// timers := b.timers
	// b.timers = make([]*timer, 0)
	// return timers
	timers := make([]*Timer, 0)
	e := b.timers.Front()
	for e != nil {
		next := e.Next()
		timer := e.Value.(*Timer)
		timers = append(timers, timer)
		b.timers.Remove(e)
		e = next
	}
	return timers
}

// runTimerTask
func (b *bucket) runTimerTask() {
	b.Lock()
	defer b.Unlock()
	e := b.timers.Front()
	for e != nil {
		next := e.Next()
		timer := e.Value.(*Timer)
		b.timers.Remove(e)
		go timer.task()
		e = next
	}
}
