package goTtrack

import (
	"runtime"
	"strings"
	"time"
)

const (
	isNormal     = iota
	isStartWatch = iota
	isTimePoint  = iota
)

//Collector Variables
var iPCchannel chan channelStruct
var stats map[string]*FuncStats
var cid int

type channelStruct struct {
	T       time.Time
	PC      uintptr
	Type    int
	WatchID int
}

//Start routine to recieve times and put them into the struct that holds all times
func init() {
	go collector()
}

//=============================================================================
// Exported Functions

/*
Measure time from start to end

Usage: defer Observer.TimeTrack(time.Now())
*/
func TimeTrack(start time.Time) {

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	//Let the Collector routine handle the rest
	iPCchannel <- channelStruct{
		T:  start,
		PC: pc,
	}
}

//Track starts a tracker that is also able to add specific TimePoints
func Track(start time.Time) int {
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	id := cid
	cid++

	//Let the Collector routine handle the rest
	iPCchannel <- channelStruct{
		T:       start,
		PC:      pc,
		Type:    isStartWatch,
		WatchID: id,
	}

	return id
}

//TimePoint adds a time point to the given id
//the user is responsible for the corretnes of the id
func TimePoint(time time.Time, id int) {
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	//Let the Collector routine handle the rest
	iPCchannel <- channelStruct{
		T:       time,
		PC:      pc,
		Type:    isTimePoint,
		WatchID: id,
	}
}

//=============================================================================
//Unexported functions

/*The Collector observes the iPCchannel and puts all recieved data into the struct*/
func collector() {
	stats = make(map[string]*FuncStats)

	iPCchannel = make(chan channelStruct, 100)

	for {
		val, ok := <-iPCchannel
		if ok {
			switch val.Type {
			case isNormal:
				addTime(&val)
				break
			case isStartWatch:
				startWatch(&val)
				break
			case isTimePoint:
				addTimePoint(&val)
				break
			}
		} else {
			break
		}
	}
}

func getFuncNameAndLine(pc uintptr) (string, int) {

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	//Get call line
	_, line := funcObj.FileLine(pc)

	name := funcObj.Name()[strings.LastIndex(funcObj.Name(), "/")+1:]

	return name, line
}

func addTime(param *channelStruct) {

	name, line := getFuncNameAndLine(param.PC)

	//Get elapsed time
	elapsed := time.Since(param.T)

	entry := stats[name]
	if entry == nil { //entry.Name should always be set when entry is valid
		entry = newFuncStats(name, line)
	} else {
		entry.Line = line
	}
	entry.Durations = append(entry.Durations, elapsed)
	stats[name] = entry
}

func startWatch(param *channelStruct) {
	name, line := getFuncNameAndLine(param.PC)

	entry := stats[name]
	if entry == nil { //entry.Name should always be set when entry is valid
		entry = newFuncStats(name, line)
	}

	t := newCTimer(line, param.T)
	entry.TimePoints[param.WatchID] = &t

	stats[name] = entry
}

func addTimePoint(param *channelStruct) {

	name, line := getFuncNameAndLine(param.PC)

	entry := stats[name]
	if entry == nil {
		panic("Uninitialized function entry")
	}

	watch := entry.TimePoints[param.WatchID]
	if watch.IsEmpty() {
		panic("Uninitialized Stopwatch")
	}

	entry.TimePoints[param.WatchID].Durations[line] = watch.Round(line, param.T)

	stats[name] = entry
}

func calcAll() {
	for k := range stats {
		stats[k].calcStats()
	}
}
