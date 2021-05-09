package common

import "time"

var ZoneOffset int32

func init() {
	_, offset := Now().Zone()
	ZoneOffset = int32(offset)
}

var NowOffset time.Duration

func Now() time.Time {
	return time.Now().Add(NowOffset)
}

func Unix() int32 {
	return int32(Now().Unix())
}

func Time(unix int32) time.Time {
	return time.Unix(int64(unix), 0)
}