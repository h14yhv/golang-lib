package clock

import "time"

func Format(value time.Time, layout string) string {
	// Success
	return value.Format(layout)
}

func Unix(value time.Time) int64 {
	// Success
	return value.Unix()
}

func UnixMilli(value time.Time) int64 {
	// success
	return int64(time.Nanosecond) * value.UnixNano() / int64(time.Millisecond)
}

func UnixMicro(value time.Time) int64 {
	// Success
	return int64(time.Nanosecond) * value.UnixNano() / int64(time.Microsecond)
}

func UnixNano(value time.Time) int64 {
	// Success
	return value.UnixNano()
}
