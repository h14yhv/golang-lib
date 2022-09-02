package clock

func (d Duration) Nanoseconds() int64 { return int64(d) }

func (d Duration) Microseconds() int64 { return int64(d) / 1e3 }

func (d Duration) Milliseconds() int64 { return int64(d) / 1e6 }

func (d Duration) Seconds() float64 {
	sec := d / Second
	nsec := d % Second
	// Success
	return float64(sec) + float64(nsec)/1e9
}

func (d Duration) Minutes() float64 {
	min := d / Minute
	nsec := d % Minute
	// Success
	return float64(min) + float64(nsec)/(60*1e9)
}

func (d Duration) Hours() float64 {
	hour := d / Hour
	nsec := d % Hour
	// Success
	return float64(hour) + float64(nsec)/(60*60*1e9)
}
