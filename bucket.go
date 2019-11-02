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
	// log.Info("add timer")
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
	defer b.Unlock()
	for _, timer := range b.timers {
		go timer.task()
	}
	b.timers = make([]*timer, 0)
}
