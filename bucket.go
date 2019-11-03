package mytimewheel

import (
	"sync"
)

type bucket struct {
	sync.Mutex
	timers []*timer
}

func newBucket() *bucket {
	return &bucket{
		timers: make([]*timer, 0),
	}
}

type timer struct {
	expire int64
	task   func()
}

func newTimer(expire int64, task func()) *timer {
	return &timer{
		expire: expire,
		task:   task,
	}
}

func (b *bucket) addTimer(t *timer) {
	b.Lock()
	defer b.Unlock()
	b.timers = append(b.timers, t)
}

func (b *bucket) getTimers() []*timer {
	b.Lock()
	defer b.Unlock()
	timers := b.timers
	b.timers = make([]*timer, 0)
	return timers
}

// runTimerTask
func (b *bucket) runTimerTask() {
	b.Lock()
	timers := b.timers
	b.timers = make([]*timer, 0)
	b.Unlock()
	for _, timer := range timers {
		go timer.task()
	}
}
