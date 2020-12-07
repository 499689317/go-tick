package ticker

import (
	"fmt"
	"testing"
)

type TestSchedule struct {
}

func (s *TestSchedule) Init() {
	tick := NewTicker(1)
	tick.AddHandle(1, -1, s)
}
func (s *TestSchedule) onTickerUpdate(nowTime int64) bool {
	fmt.Println("update", nowTime)
	return true
}
func (s *TestSchedule) onTickerExit(nowTime int64) bool {
	fmt.Println("exit", nowTime)
	return true
}

func TestNewTicker(t *testing.T) {
	ts := new(TestSchedule)
	ts.Init()
	select {}
}
