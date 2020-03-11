package goTtrack

import (
	"fmt"
	"sort"
	"strconv"
)

/*GetFuncStats gets a specific Statistic entry */
func GetFuncStats(name string) *FuncStats {
	entry := stats[name]
	if entry != nil {
		entry.calcStats()
	}
	return entry
}

/*GetStats gets all current Statistic entrys */
func GetStats() map[string]*FuncStats {
	calcAll()
	return stats
}

func GetStatsPrint(printRaw bool) (ret string) {
	ret += fmt.Sprintln("====== Results =====")

	//TODO Sort map by keys

	for _, v := range stats {

		ret += fmt.Sprintf("%v:%v\n", v.Name, strconv.Itoa(v.Line))

		//Func execs
		if len(v.Durations) > 0 {
			ret += fmt.Sprintln("Full runs: ", v.Stats.Count, "-> Median:", v.Stats.Med, "Min:", v.Stats.Min, "Max:", v.Stats.Max)
			if printRaw {
				ret += fmt.Sprintln("Durations:", v.Durations)
			}
		}

		//func timepoints
		if len(v.TimePoints) > 0 {
			var keys []int
			for k := range v.TimePoints {
				keys = append(keys, k)
			}
			sort.Ints(keys)

			ret += fmt.Sprintln("Timepoints:")
			for _, l := range keys {
				ret += fmt.Sprintln("ID", v.TimePoints[l].Name, ":", v.TimePoints[l].Durations)
				ret += fmt.Sprintln("ID", l, ":", v.TimePoints[l].Stats)
			}
		}

		ret += fmt.Sprintln()
	}
	return ret
}

//PrintAllStats waits until the buffer of the ipc channel is empty
func GetCalcStatsPrint(printValues bool) (ret string) {
	for len(iPCchannel) != 0 {
	}
	calcAll()
	return GetStatsPrint(printValues)
}
