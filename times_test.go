package goTtrack

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	id := Track(time.Now())

	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 10)

		TimePoint(time.Now(), id)

		dummy()
	}

	fmt.Println(GetCalcStatsPrint(true))
}

func dummy() {
	defer TimeTrack(time.Now())

	time.Sleep(time.Millisecond * 20)
}
