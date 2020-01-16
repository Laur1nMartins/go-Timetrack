package goTtrack

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strconv"
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
var stats map[string]FuncStats
var isInit bool

/*FuncStats holds all currently available stats of the function*/
type FuncStats struct {
	Name string //Name of the function
	Line int    //Line that returned from the tracked function

	Durations []time.Duration //Execution times

	TimePoints map[int]CTimer

	Count  int //Count how many times the function was called and returned
	isInit bool

	Average time.Duration //Average Execution time
	Min     time.Duration //Minimal Execution time
	Max     time.Duration //Maximal Execution time
}

type channelStruct struct {
	T       time.Time
	PC      uintptr
	Type    int
	WatchID int
}

//=============================================================================
// Struct Functions

func (f *FuncStats) initFuncStats() {
	if f.isInit == false {
		f.Durations = make([]time.Duration, 0)
		f.TimePoints = make(map[int]CTimer)
		f.isInit = true
	}
}

/*Calculates the Execution times of a given Function (Average, Min, Max)*/
func (f *FuncStats) calcStats() {

	var all time.Duration

	f.Count = len(f.Durations)

	if f.Count == 0 {
		return
	}

	for i, d := range f.Durations {
		all += d
		if i == 0 {
			f.Min = d
		}
		if d > f.Max {
			f.Max = d
		} else {
			if d < f.Min {
				f.Min = d
			}
		}
	}

	f.Average = time.Duration(int64(all) / int64(f.Count))
}

func calcAll() {
	for k := range stats {
		tmp := stats[k]
		tmp.calcStats()
	}
}

/*The Collector observes the iPCchannel and puts all recieved data into the struct*/
func collector() {
	stats = make(map[string]FuncStats)

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

//=============================================================================
// Exported Functions
// StartCollector needs to be started at first

/*StartCollector starts the collector routine. */
func Start() {
	if isInit {
		log.Println("Collector already running!")
		return
	}
	go collector()
	isInit = true
}

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

func StartWatch(start time.Time, id int) {
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	//Let the Collector routine handle the rest
	iPCchannel <- channelStruct{
		T:       start,
		PC:      pc,
		Type:    isStartWatch,
		WatchID: id,
	}
}

//TimePoint adds a time point to the
func TimePoint(start time.Time, id int) {
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	//Let the Collector routine handle the rest
	iPCchannel <- channelStruct{
		T:       start,
		PC:      pc,
		Type:    isTimePoint,
		WatchID: id,
	}
}

/*GetFuncStats gets a specific Statistic entry */
func GetFuncStats(name string) FuncStats {
	tmp := stats[name]
	if tmp.isInit {
		tmp.calcStats()
	}
	return tmp
}

/*GetStats gets all current Statistic entrys */
func GetStats() map[string]FuncStats {
	calcAll()
	return stats
}

func GetStatsPrint() (ret string) {
	ret += fmt.Sprintln("====== Results =====")
	for _, v := range stats {
		ret += fmt.Sprintln(v.Name+":"+strconv.Itoa(v.Line), v.Count, "Avg:", v.Average, "Min:", v.Min, "Max:", v.Max)
		ret += fmt.Sprintln("Durations:", v.Durations)
		var keys []int
		for k := range v.TimePoints {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		ret += fmt.Sprintln("Timepoints:")
		for _, l := range keys {
			ret += fmt.Sprintln("ID", l, ":", v.TimePoints[l].Durations)
		}
		ret += fmt.Sprintln()
	}
	return ret
}

//PrintAllStats waits until the buffer of the ipc channel is empty
func GetAllStatsPrint() (ret string) {
	for len(iPCchannel) != 0 {
	}
	calcAll()
	return GetStatsPrint()
}

//=============================================================================
// Helper Functions

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
	if !entry.isInit { //entry.Name should always be set when entry is valid
		entry.initFuncStats()
		entry.Line = line
		entry.Name = name
	} else {
		entry.Line = line
	}
	entry.Durations = append(entry.Durations, elapsed)
	stats[name] = entry
}

func startWatch(param *channelStruct) {
	name, line := getFuncNameAndLine(param.PC)

	entry := stats[name]
	if !entry.isInit { //entry.Name should always be set when entry is valid
		entry.initFuncStats()
		entry.Line = line
		entry.Name = name
	}

	entry.TimePoints[param.WatchID] = StartCTimer(line, param.T)

	stats[name] = entry
}

func addTimePoint(param *channelStruct) {

	name, line := getFuncNameAndLine(param.PC)

	entry := stats[name]

	if !entry.isInit {
		panic("Uninitialized function entry")
	}

	watch := entry.TimePoints[param.WatchID]

	if watch.IsEmpty() {
		panic("Uninitialized Stopwatch")
	}

	entry.TimePoints[param.WatchID].Durations[line] = watch.Round(line, param.T)

	stats[name] = entry
}
