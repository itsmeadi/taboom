package rate

import (
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	rate := 10
	noOfRequest := 500
	timeToProcess := 100 * time.Millisecond

	r := InitRateLimiter(timeToProcess, rate)

	t1 := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < noOfRequest; i++ {
		wg.Add(1)
		go func(*sync.WaitGroup) {
			r.Wait()
			time.Sleep(800 * time.Millisecond)

			wg.Done()
		}(&wg)

	}
	wg.Wait()
	if (timeToProcess * time.Duration(noOfRequest/rate)) > time.Now().Sub(t1) {
		t.Fatalf("Expected=%+v, GOT=%+v", (timeToProcess * time.Duration(noOfRequest/rate)), time.Now().Sub(t1))
	}
}
