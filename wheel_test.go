package mytimewheel

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	w := NewWheel(1, 1024)
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
	w := NewWheel(1, 1024)
	if w == nil {

	}
	var wg sync.WaitGroup
	taskCnt := 100000
	wg.Add(taskCnt)
	for i := 0; i < taskCnt; i++ {
		d := rand.Int63n(10000) + 1
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
	file, _ := os.Create("cpu_profile")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()
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
}

func TestCpuUseIngo(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

func TestTicker(t *testing.T) {
	w := NewWheel(1, 1024)
	ticker := w.NewTicker(1000 * time.Millisecond)
	for {
		<-ticker.C
		log.Println("run ticker trigger")
	}
}

func TestManyIngoTicker(t *testing.T) {
	for i := 0; i < 100000; i++ {
		d := rand.Int63n(10000) + 1
		ticker := time.NewTicker(time.Duration(d) * time.Millisecond)
		go func(tc *time.Ticker) {
			for {
				<-ticker.C
			}
		}(ticker)
	}
	time.Sleep(10000 * time.Second)
}

func TestManyMyTicker(t *testing.T) {
	w := NewWheel(1, 1024)
	for i := 0; i < 100000; i++ {
		d := rand.Int63n(10000) + 1
		ticker := w.NewTicker(time.Duration(d) * time.Millisecond)
		go func(tc *MyTicker) {
			for {
				<-tc.C
			}
		}(ticker)
	}
	if true {

	}
	time.Sleep(10000 * time.Second)
}
func TestManySchedule(t *testing.T) {
	w := NewWheel(1, 1024)

	for i := 0; i < 100000; i++ {
		d := rand.Int63n(10000) + 1
		w.Schedue(time.Duration(d)*time.Millisecond, func() {})
	}
	if true {

	}
	time.Sleep(10000 * time.Second)
}
func TestSchedule(t *testing.T) {
	w := NewWheel(1, 1024)
	var cnt int32
	w.Schedue(500*time.Millisecond, func() {
		log.Println("ticker...")
		atomic.AddInt32(&cnt, 1)
	})
	time.Sleep(5 * time.Second)
	if cnt != 10 {
		t.Errorf("schedule cnt err cnt:%d", cnt)
	}
}
func TestTimesSchedule(t *testing.T) {
	w := NewWheel(1, 1024)
	var cnt int
	w.SchedueWithTimes(time.Second, 3, func() {
		log.Println("ticker...")
		cnt++
	})
	time.Sleep(6 * time.Second)
	if cnt != 3 {
		t.Errorf("run times fail:%d", cnt)
	}

}
func TestIngoSigleTicker(t *testing.T) {
	er := time.NewTicker(time.Second)
	go func() {
		time.Sleep(2 * time.Second)
		er.Stop()
	}()
	for {
		<-er.C
		fmt.Println("triget")
		time.Sleep(2 * time.Second)
	}

}

func TestGetBucketSize(t *testing.T) {
	size := getBucketSize(1025)
	t.Logf("ret:%d", size)
	if size != 2048 {
		t.Errorf("get bucket size:%d err", size)
	}
	size = getBucketSize(1024)
	if size != 1024 {
		t.Errorf("get bucket size 1024 err real:%d", size)
	}
}
