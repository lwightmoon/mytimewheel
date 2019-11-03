package mytimewheel

import (
	"math"
	"sync/atomic"
	"time"
	"unsafe"
)

type Wheel struct {
	tick     int64
	size     int64
	interval int64
	ticker   *time.Ticker
	child    *Wheel
	parent   unsafe.Pointer
	// curTime  int64
	pos     int64 //此刻的时间轮
	buckets []*bucket
}

func NewDefaultWheel() *Wheel {
	return NewWheel(1, 1024)
}

func NewWheel(tickMs, size int64) *Wheel {
	size = getBucketSize(size)
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	ticker := time.NewTicker(time.Duration(tickMs) * time.Millisecond)
	wheel := &Wheel{
		tick:     tickMs,
		size:     size,
		interval: tickMs * size,
		ticker:   ticker,
		buckets:  buckets,
	}
	wheel.run()
	return wheel
}

func (w *Wheel) add(t *timer) bool {
	// curtime := atomic.LoadInt64(&w.curTime)
	curtime := time.Now().UnixNano() / int64(time.Millisecond)
	if t.expire < curtime+w.tick {
		return false
	} else if t.expire < curtime+w.interval {
		index := (t.expire - curtime) / w.tick
		if index > 0 {
			index--
		}
		// index = (w.pos + index) % w.size
		index = (w.pos + index) & (w.size - 1)
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
		ret := (*Wheel)(parentWheel).add(t)
		return ret
	}
}

func (w *Wheel) run() {
	go func() {
		for {
			<-w.ticker.C
			bucket := w.buckets[w.pos]
			// w.pos = (w.pos + 1) % w.size
			w.pos = (w.pos + 1) & (w.size - 1)
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

func (w *Wheel) addOrRun(t *timer) {
	if !w.add(t) {
		go t.task()
	}
}

func (w *Wheel) AfterFunc(d time.Duration, task func()) {
	expire := time.Now().Add(d).UnixNano() / int64(time.Millisecond)
	t := newTimer(expire, task)
	w.addOrRun(t)
}

type MyTicker struct {
	stopflag int32
	C        chan struct{}
}

func (mt *MyTicker) GetC() <-chan struct{} {
	return mt.C
}

func (mt *MyTicker) Stop() {
	atomic.StoreInt32(&mt.stopflag, 1)
}

func (w *Wheel) NewTicker(d time.Duration) *MyTicker {
	myTicker := &MyTicker{
		stopflag: 0,
		C:        make(chan struct{}, 1),
	}
	var t *timer
	t = &timer{
		expire: time.Now().Add(d).UnixNano() / int64(time.Millisecond),
		task: func() {
			if atomic.LoadInt32(&myTicker.stopflag) != 1 {
				select {
				case myTicker.C <- struct{}{}:
				default:
				}
				expire := t.expire + int64(d/time.Millisecond)
				t.expire = expire
				w.addOrRun(t)
			} else {
				close(myTicker.C)
			}
		},
	}
	w.addOrRun(t)
	return myTicker
}

func (w *Wheel) Schedue(d time.Duration, f func()) {
	var t *timer
	t = &timer{
		expire: time.Now().Add(d).UnixNano() / int64(time.Millisecond),
		task: func() {
			expire := t.expire + int64(d/time.Millisecond)
			t.expire = expire
			w.addOrRun(t)
			f()
		},
	}
	w.addOrRun(t)
}

func (w *Wheel) SchedueWithTimes(d time.Duration, times int32, f func()) {
	var t *timer
	var cnt int32
	if times < 1 {
		w.Schedue(d, f)
		return
	}
	t = &timer{
		expire: time.Now().Add(d).UnixNano() / int64(time.Millisecond),
		task: func() {
			if atomic.LoadInt32(&cnt) < times {
				expire := t.expire + int64(d/time.Millisecond)
				t.expire = expire
				w.addOrRun(t)
				f()
				atomic.AddInt32(&cnt, 1)
			}
		},
	}
	w.addOrRun(t)
}

// copy java hashmap tableSizeFor 取>num的最小2的n次幂
func getBucketSize(num int64) int64 {
	num = num - 1
	num |= num >> 1
	num |= num >> 2
	num |= num >> 4
	num |= num >> 8
	num |= num >> 16
	num |= num >> 32
	if num < 0 {
		return 1
	}
	if num > math.MaxInt64 {
		return math.MaxInt64
	}
	return num + 1
}
