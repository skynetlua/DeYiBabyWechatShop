package common

import (
	"sync/atomic"
	"time"
)

func BeginWrite(writeCounter *int32) {
	if *writeCounter == 0 {
		atomic.AddInt32(writeCounter, 1)
	}else{
		if *writeCounter != 0 {
			waitTime := time.Now().Unix()+3
			waitCount := 0
			for *writeCounter != 0 {
				waitCount++
				if waitCount%10000 == 0 {
					if time.Now().Unix() > waitTime {
						panic("BeginWrite is lock")
					}
				}
			}
		}
		atomic.AddInt32(writeCounter, 1)
	}
}

func EndWrite(writeCounter *int32) {
	if *writeCounter == 0 {
		panic("EndWrite is finish")
	}else{
		atomic.AddInt32(writeCounter, -1)
	}
}