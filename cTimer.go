package goTtrack

import (
	"time"
)

type CTimer struct {
	StartT    time.Time
	StartLine int
	Durations map[int][]time.Duration //int == line
}

func (c CTimer) IsEmpty() bool {
	return c.StartLine == 0 && c.StartT.IsZero()
}

func (c CTimer) Round(line int, end time.Time) []time.Duration {
	e := c.Durations[line]
	if e == nil {
		c.Durations[line] = make([]time.Duration, 0)
	}
	return append(c.Durations[line], end.Sub(c.StartT))
}

func StartCTimer(line int, startTime time.Time) (c CTimer) {
	return CTimer{
		StartT:    startTime,
		StartLine: line,
		Durations: make(map[int][]time.Duration, 0),
	}
}
