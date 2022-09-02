package clock

import (
	"time"
)

func Now(timezone string) (time.Time, error) {
	if timezone == "" {
		timezone = UTC
	}
	now := time.Now()
	// Success
	return InTimezone(now, timezone)
}

func LoadLocation(tz string) (*time.Location, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	// Success
	return loc, nil
}

func New(year, month, day, hour, minute, second, nanosecond int, timezone string) (time.Time, error) {
	loc, err := LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	// Success
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, loc), nil
}

func Replace(value time.Time, year, month, day, hour, minute, second, nanosecond int, timezone string) (time.Time, error) {
	nYear := value.Year()
	nMonth := int(value.Month())
	nDay := value.Day()
	nHour := value.Hour()
	nMinute := value.Minute()
	nSecond := value.Second()
	nNanosecond := value.Nanosecond()
	nTimezone := value.Location()
	if year != Null {
		nYear = year
	}
	if month != Null {
		nMonth = month
	}
	if day != Null {
		nDay = day
	}
	if hour != Null {
		nHour = hour
	}
	if minute != Null {
		nMinute = minute
	}
	if second != Null {
		nSecond = second
	}
	if nanosecond != Null {
		nNanosecond = nanosecond
	}
	if timezone != "" {
		tz, err := LoadLocation(timezone)
		if err != nil {
			return time.Time{}, err
		}
		nTimezone = tz
	}
	// Success
	return time.Date(nYear, time.Month(nMonth), nDay, nHour, nMinute, nSecond, nNanosecond, nTimezone), nil
}

func ReplaceTimezone(value time.Time, timezone string) (time.Time, error) {
	// Success
	return Replace(value, Null, Null, Null, Null, Null, Null, Null, timezone)
}

func ReplaceDate(value time.Time, year, month, day int) (time.Time, error) {
	// Success
	return Replace(value, year, month, day, Null, Null, Null, Null, "")
}

func ReplaceTime(value time.Time, hour, minute, second, nanosecond int) (time.Time, error) {
	// Success
	return Replace(value, Null, Null, Null, hour, minute, second, nanosecond, "")
}

func InTimezone(value time.Time, timezone string) (time.Time, error) {
	loc, err := LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	// Success
	return value.In(loc), nil
}

func Sleep(duration Duration) {
	// Success
	time.Sleep(time.Duration(duration))
}
