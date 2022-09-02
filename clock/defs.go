package clock

import "time"

const (
	// Format RFC 3339
	FormatRFC3339Z       = "2006-01-02T15:04:05Z07:00"
	FormatRFC3339MilliZ  = "2006-01-02T15:04:05.999Z07:00"
	FormatRFC3339MicroZ  = "2006-01-02T15:04:05.999999Z07:00"
	FormatRFC3339NanoZ   = "2006-01-02T15:04:05.999999999Z07:00"
	FormatRFC3339        = "2006-01-02T15:04:05"
	FormatRFC3339Milli   = "2006-01-02T15:04:05.999"
	FormatRFC3339Micro   = "2006-01-02T15:04:05.999999"
	FormatRFC3339Nano    = "2006-01-02T15:04:05.999999999"
	FormatRFC3339Date    = "2006-01-02"
	FormatRFC3339CZ      = "2006-01-02T15:04:05Z07:00"
	FormatRFC3339CMilliZ = "2006-01-02 15:04:05.999Z07:00"
	FormatRFC3339CMicroZ = "2006-01-02 15:04:05.999999Z07:00"
	FormatRFC3339CNanoZ  = "2006-01-02 15:04:05.999999999Z07:00"
	FormatRFC3339C       = "2006-01-02 15:04:05"
	FormatRFC3339CMilli  = "2006-01-02 15:04:05.999"
	FormatRFC3339CMicro  = "2006-01-02 15:04:05.999999"
	FormatRFC3339CNano   = "2006-01-02 15:04:05.999999999"
	FormatRFC3339HMZ     = "2006-01-02T15:04Z07:00"
	FormatRFC3339HM      = "2006-01-02T15:04"
	// Format Common
	FormatHuman       = "02/01/2006 15:04:05"
	FormatHumanZ      = "02/01/2006 15:04:05Z"
	FormatHumanMilli  = "02/01/2006 15:04:05.000"
	FormatHumanMilliZ = "02/01/2006 15:04:05.000Z"
	FormatHumanDate   = "02/01/2006"
	FormatHumanTime   = "15:04:05"
	FormatANSIC       = "Mon Jan _2 15:04:05 2006"
	FormatUnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	FormatRubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	FormatKitchen     = "3:04PM"
	FormatIndexYear   = "2006"
	FormatIndexMonth  = "200601"
	FormatIndexDate   = "20060102"
	// Format RFC 822
	FormatRFC822  = "02 Jan 06 15:04 MST"
	FormatRFC822Z = "02 Jan 06 15:04 -0700"
	// Format RFC 850
	FormatRFC850 = "Monday, 02-Jan-06 15:04:05 MST"
	// Format RFC 1123
	FormatRFC1123  = "Mon, 02 Jan 2006 15:04:05 MST"
	FormatRFC1123Z = "Mon, 02 Jan 2006 15:04:05 -0700"
	// Format Stamps
	FormatStamp      = "Jan _2 15:04:05"
	FormatStampMilli = "Jan _2 15:04:05.000"
	FormatStampMicro = "Jan _2 15:04:05.000000"
	FormatStampNano  = "Jan _2 15:04:05.000000000"
	// Timezone
	UTC     = "UTC"
	Local   = "Asia/Ho_Chi_Minh"
	VietNam = "Asia/Ho_Chi_Minh"
	// Common
	Null = -1

	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
	Week                 = 7 * Day

	ParseTimeError = "parse time error"
)

type (
	Duration time.Duration
)

var (
	GroupRFC3339 = []string{
		FormatRFC3339Z,
		FormatRFC3339,
		FormatRFC3339MilliZ,
		FormatRFC3339Milli,
		FormatRFC3339MicroZ,
		FormatRFC3339Micro,
		FormatRFC3339NanoZ,
		FormatRFC3339Nano,
		FormatRFC3339Date,
		FormatRFC3339CZ,
		FormatRFC3339CMilliZ,
		FormatRFC3339CMicroZ,
		FormatRFC3339CNanoZ,
		FormatRFC3339C,
		FormatRFC3339CMilli,
		FormatRFC3339CMicro,
		FormatRFC3339CNano,
		FormatRFC3339HMZ,
		FormatRFC3339HM,
	}

	GroupRFC822 = []string{
		FormatRFC822Z,
		FormatRFC822,
	}

	GroupRFC850 = []string{
		FormatRFC850,
	}

	GroupRFC1123 = []string{
		FormatRFC1123Z,
		FormatRFC1123,
	}

	GroupStamps = []string{
		FormatStamp,
		FormatStampMilli,
		FormatStampMicro,
		FormatStampNano,
	}

	GroupCommon = []string{
		FormatHuman,
		FormatHumanZ,
		FormatHumanMilli,
		FormatHumanMilliZ,
		FormatHumanDate,
		FormatHumanTime,
		FormatANSIC,
		FormatUnixDate,
		FormatRubyDate,
		FormatKitchen,
	}

	GroupIndex = []string{
		FormatIndexYear,
		FormatIndexMonth,
		FormatIndexDate,
	}

	Group = [][]string{
		GroupRFC3339,
		GroupRFC1123,
		GroupRFC850,
		GroupRFC822,
		GroupStamps,
		GroupCommon,
		GroupIndex,
	}
)
