package rate

import (
	"sync"
	"time"
)

type RateLimit struct {
	Time     time.Duration
	Rate     int
	channel  chan int
	prevTime time.Time
	lock sync.Mutex

}

func InitRateLimiter(timeLag time.Duration, rate int) RateLimit {

	ch := make(chan int, rate)
	instance := RateLimit{
		Time:     timeLag,
		Rate:     rate,
		channel:  ch,
		prevTime: time.Now(),
	}
	instance.fill()
	return instance
}
func (r RateLimit) fill() {
	for i := 0; i < r.Rate; i++ {
		r.channel <- i
	}
}

func (r *RateLimit) Wait() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if len(r.channel) == 0 {
		r.fill()
		diff := time.Now().Sub(r.prevTime)
		time.Sleep(r.Time - diff)
	}
	<-r.channel
	r.prevTime = time.Now()
}
