package goTtrack

import (
	"time"
)

//FuncStats holds all currently available stats of the function
//The name is assumed to be a unique for every function
type FuncStats struct {
	Name string //Name of the function
	Line int    //Line where initial call to track is located

	//Full execution Variables
	Durations []time.Duration //Execution times of full executions
	Stats     meta

	//Mid function time points
	TimePoints map[int]*CTimer //int -> line inside function
}

type CTimer struct {
	StartT    time.Time
	StartLine int
	Durations map[int][]time.Duration //int == line
	Stats     map[int]meta
}

type meta struct {
	Count int
	Min   time.Duration
	Max   time.Duration
	Med   time.Duration
	Avg   time.Duration
}

//=============================================================================
//

//
func newFuncStats(name string, line int) *FuncStats {
	return &FuncStats{
		Name:       name,
		Line:       line,
		Durations:  make([]time.Duration, 0),
		TimePoints: make(map[int]*CTimer),
	}
}

//Calculates the Execution times of a given Function (Med, Min, Max)
func (f *FuncStats) calcStats() {

	f.Stats.Count = len(f.Durations)

	//Iterate over durations
	if f.Stats.Count > 0 {
		f.Stats.Min = f.Durations[0]
		f.Stats.Max = f.Durations[0]
		for _, d := range f.Durations {
			if d > f.Stats.Max {
				f.Stats.Max = d
			} else {
				if d < f.Stats.Min {
					f.Stats.Min = d
				}
			}
		}
		f.Stats.Med = f.Durations[len(f.Durations)/2]
	}

	//Iterate over each individual TimePoint

}

//
func newCTimer(line int, startTime time.Time) (c CTimer) {
	return CTimer{
		StartT:    startTime,
		StartLine: line,
		Durations: make(map[int][]time.Duration, 0),
	}
}

func (c *CTimer) calc() {

}

func (c *CTimer) IsEmpty() bool {
	return c.StartLine == 0 && c.StartT.IsZero()
}

func (c CTimer) Round(line int, end time.Time) []time.Duration {
	e := c.Durations[line]
	if e == nil {
		c.Durations[line] = make([]time.Duration, 0)
	}
	return append(c.Durations[line], end.Sub(c.StartT))
}
