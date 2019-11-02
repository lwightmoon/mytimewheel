package mytimewheel

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type wheel struct {
	tick     int64
	size     int64
	interval int64
	ticker   *time.Ticker
	child    *wheel
	parent   unsafe.Pointer
	// curTime  int64
	pos     int64 //此刻的时间轮
	buckets []*bucket
}

func NewWheel(tickMs, size int64) *wheel {
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	ticker := time.NewTicker(time.Duration(tickMs) * time.Millisecond)
	wheel := &wheel{
		tick:     tickMs,
		size:     size,
		interval: tickMs * size,
		ticker:   ticker,
		buckets:  buckets,
	}
	wheel.run()
	return wheel
}

func (w *wheel) add(t *timer) bool {
	// curtime := atomic.LoadInt64(&w.curTime)
	curtime := time.Now().UnixNano() / int64(time.Millisecond)
	if t.expire < curtime+w.tick {
		return false
	} else if t.expire < curtime+w.interval {
		index := (t.expire - curtime) / w.tick
		if index > 0 {
			index--
		}
		index = (w.pos + index) % w.size
		b := w.buckets[index]
		b.addTimer(t)
		return true
	} else {

		parentWheel := atomic.LoadPointer(&w.parent)
		if parentWheel == nil {
			parent := NewWheel(w.interval, w.size)
			parent.child = w
			atomic.CompareAndSwapPointer(
				&w.parent,
				nil,
				unsafe.Pointer(parent),
			)
			parentWheel = atomic.LoadPointer(&w.parent)
		}
		ret := (*wheel)(parentWheel).add(t)
		return ret
	}
}

func (w *wheel) run() {
	go func() {
		for {
			<-w.ticker.C
			bucket := w.buckets[w.pos]
			w.pos = (w.pos + 1) % w.size
			// w.pos = (w.pos + 1) & (w.size - 1)
			if w.child == nil { //最底层
				bucket.runTimerTask()
			} else { //上层
				timers := bucket.getTimers()
				if len(timers) > 0 {
					for _, timer := range timers {
						if !w.child.add(timer) {
							go timer.task()
						}
					}
				} else {
				}
			}
		}
	}()
}

func (w *wheel) AfterFunc(d time.Duration, task func()) {
	expire := time.Now().Add(d).UnixNano() / int64(time.Millisecond)
	t := newTimer(expire, task)
	if !w.add(t) {
		task()
	}
}
