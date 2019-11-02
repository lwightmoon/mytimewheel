package mytimewheel

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	w := NewWheel(1, 1000)
	now := time.Now().Unix()
	w.AfterFunc(2*time.Second, func() {
		end := time.Now().Unix()
		t.Logf("cost1:%d\n", end-now)
		cost := end - now
		if cost != 2 {
			t.Errorf("not equal cost1:%d", cost)
		}
	})
	time.Sleep(3 * time.Second)
	now2 := time.Now().Unix()
	w.AfterFunc(2*time.Second, func() {
		end := time.Now().Unix()
		t.Logf("cost2:%d\n", end-now2)
		cost := end - now2
		if cost != 2 {
			t.Errorf("not equal cost2:%d", cost)
		}
	})
	time.Sleep(3 * time.Second)
}

func TestTask(t *testing.T) {
	var cnt int64
	runtime.GOMAXPROCS(runtime.NumCPU())
	w := NewWheel(1, 1000)
	if w == nil {

	}
	var wg sync.WaitGroup
	taskCnt := 100000
	wg.Add(taskCnt)
	for i := 0; i < taskCnt; i++ {
		d := rand.Int63n(10000)
		// start := time.Now().UnixNano() / int64(time.Millisecond)
		w.AfterFunc(time.Duration(d)*time.Millisecond, func() {
			defer wg.Done()
			atomic.AddInt64(&cnt, 1)
		})
		//
		// time.AfterFunc(time.Duration(d)*time.Millisecond, func() {
		//defer wg.Done()
		// end := time.Now().UnixNano() / int64(time.Millisecond)
		// cost := end - start
		// if (cost - d) > 10 {
		// 	t.Errorf("after err real:%d,expect:%d", cost, d)
		// }
		//})
	}
	wg.Wait()
	if cnt != int64(taskCnt) {
		t.Error("task not all exec")
	}
}

func TestCpuUseMy(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	w := NewWheel(1, 1024)
	if w == nil {
		t.Log("")
	}
	var wg sync.WaitGroup
	taskCnt := 100000
	wg.Add(taskCnt)
	for i := 0; i < taskCnt; i++ {
		d := rand.Int63n(10000)
		w.AfterFunc(time.Duration(d)*time.Millisecond, func() {
			defer wg.Done()
		})
	}
	wg.Wait()
	//
}

func TestCpuUseIngo(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// w := newWheel(1, 1000)
	// if w == nil {

	// }
	var wg sync.WaitGroup
	taskCnt := 100000
	wg.Add(taskCnt)
	for i := 0; i < taskCnt; i++ {
		d := rand.Int63n(10000)
		time.AfterFunc(time.Duration(d)*time.Millisecond, func() {
			defer wg.Done()
		})
	}
	wg.Wait()
}

func TestSleep(t *testing.T) {
	size := 1024
	a := 2
	for {
		time.Sleep(time.Millisecond)
		b := a % size
		if b == 0 {

		}
	}
}
