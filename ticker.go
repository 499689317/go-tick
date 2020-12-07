package ticker

import (
	"sync"
	"time"
)

// Handle 对外接口
type Handle interface {
	OnTickerUpdate(nowTime int64) bool
	OnTickerExit(nowTime int64) bool
}

// TickHandle tick handle
type TickHandle struct {
	delayTime int64
	longTime  int64
	dealTime  int64
	handle    Handle
}

// Ticker 定时器
type Ticker struct {
	sync.Mutex
	handles []TickHandle
	dels    []TickHandle
	adds    []TickHandle
}

// NewTicker interval
func NewTicker(delay time.Duration) *Ticker {
	t := new(Ticker)
	timer := time.NewTimer(delay * time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
				now := time.Now().Unix()
				t.onTimerHandle(now)
				timer.Reset(time.Second)
			}
		}
	}()
	return t
}

// AddHandle add
func (t *Ticker) AddHandle(delayTime, longTime int64, h Handle) {
	if longTime <= 0 {
		longTime = 0xFFFFFFFF
	}
	dt := time.Now().Unix() + delayTime
	t.Lock()
	defer t.Unlock()
	t.adds = append(t.adds, TickHandle{delayTime: delayTime, longTime: longTime, dealTime: dt, handle: h})
}

// DelHandle del
func (t *Ticker) DelHandle(h Handle) {
	t.Lock()
	defer t.Unlock()
	t.dels = append(t.dels, TickHandle{handle: h})
}

func (t *Ticker) updateHandle() {
	t.Lock()
	defer t.Unlock()
	t.handles = append(t.handles, t.adds...)
	t.adds = make([]TickHandle, 0)
	for i := 0; i < len(t.dels); i++ {
		for j := 0; j < len(t.handles); j++ {
			if t.dels[i].handle == t.handles[j].handle {
				t.handles = append(t.handles[:j], t.handles[j+1:]...)
			}
		}
	}
}
func (t *Ticker) onTimerHandle(nowTime int64) {
	t.updateHandle()
	isDel := false
	for i := 0; i < len(t.handles); i++ {
		p := &t.handles[i]
		if nowTime > p.dealTime {
			if p.handle.OnTickerUpdate(nowTime) {
				p.dealTime += p.delayTime
				p.longTime -= p.delayTime
				if p.longTime <= 0 {
					isDel = true
				}
			} else {
				isDel = true
			}
		}
		if isDel {
			p.handle.OnTickerExit(nowTime)
			t.handles = append(t.handles[0:i], t.handles[i+1:]...)
			i--
		}
	}
}
