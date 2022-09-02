package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/h14yhv/golang-lib/clock"
)

var (
	// Map for the various codes of colors
	colors map[LogLevel]string
	// Map from format's placeholders to printf verbs
	phfs map[string]string
	// Contains color strings for stdout
	logNo uint64
	// Default format of log message
	defFmt = "[%[2]s][%[3]s][%[4]s:%[5]d][%.3[6]s] ▶ %[7]s"
	// Default format of time
	defTimeFmt = clock.FormatRFC3339
)

type LogLevel int

const (
	Black = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	CriticalLevel LogLevel = iota + 1
	ErrorLevel
	WarningLevel
	NoticeLevel
	InfoLevel
	DebugLevel
)

type Worker struct {
	Minion     *log.Logger
	Color      bool
	format     string
	timeFormat string
	level      LogLevel
}

type Info struct {
	Id       uint64
	Time     string
	Module   string
	Level    LogLevel
	Line     int
	Filename string
	Message  string
}

type logger struct {
	Module string
	worker *Worker
}

func init() {
	initColors()
	initFormatPlaceholders()
}

func (r *Info) Output(format string) string {
	msg := fmt.Sprintf(format,
		r.Id,               // %[1] // %{id}
		r.Time,             // %[2] // %{time[:fmt]}
		r.Module,           // %[3] // %{module}
		r.Filename,         // %[4] // %{filename}
		r.Line,             // %[5] // %{line}
		r.logLevelString(), // %[6] // %{level}
		r.Message,          // %[7] // %{message}
	)
	if i := strings.LastIndex(msg, "%!(EXTRA"); i != -1 {
		return msg[:i]
	}
	return msg
}

func parseFormat(format string) (msgfmt, timefmt string) {
	if len(format) < 10 {
		return defFmt, defTimeFmt
	}
	timefmt = defTimeFmt
	idx := strings.IndexRune(format, '%')
	for idx != -1 {
		msgfmt += format[:idx]
		format = format[idx:]
		if len(format) > 2 {
			if format[1] == '{' {
				if jdx := strings.IndexRune(format, '}'); jdx != -1 {
					idx = strings.Index(format[1:], "%{")
					if idx != -1 && idx < jdx {
						msgfmt += "%%"
						format = format[1:]
						continue
					}
					verb, arg := ph2verb(format[:jdx+1])
					msgfmt += verb
					if verb == `%[2]s` && arg != "" {
						timefmt = arg
					}
					format = format[jdx+1:]
				} else {
					format = format[1:]
				}
			} else {
				msgfmt += "%%"
				format = format[1:]
			}
		}
		idx = strings.IndexRune(format, '%')
	}
	msgfmt += format
	return
}

func ph2verb(ph string) (verb string, arg string) {
	n := len(ph)
	if n < 4 {
		return ``, ``
	}
	if ph[0] != '%' || ph[1] != '{' || ph[n-1] != '}' {
		return ``, ``
	}
	idx := strings.IndexRune(ph, ':')
	if idx == -1 {
		return phfs[ph], ``
	}
	verb = phfs[ph[:idx]+"}"]
	arg = ph[idx+1 : n-1]
	return
}

func NewWorker(prefix string, flag int, color bool, out io.Writer) *Worker {
	return &Worker{Minion: log.New(out, prefix, flag), Color: color, format: defFmt, timeFormat: defTimeFmt}
}

func SetDefaultFormat(format string) {
	defFmt, defTimeFmt = parseFormat(format)
}

func (w *Worker) SetFormat(format string) {
	w.format, w.timeFormat = parseFormat(format)
}

func (l *logger) SetFormat(format string) {
	l.worker.SetFormat(format)
}

func (w *Worker) SetLogLevel(level LogLevel) {
	w.level = level
}

func (l *logger) SetLogLevel(level LogLevel) {
	l.worker.level = level
}

func (w *Worker) Log(level LogLevel, calldepth int, info *Info) error {

	if w.level < level {
		return nil
	}
	if w.Color {
		buf := &bytes.Buffer{}
		buf.Write([]byte(colors[level]))
		buf.Write([]byte(info.Output(w.format)))
		buf.Write([]byte("\033[0m"))
		return w.Minion.Output(calldepth+1, buf.String())
	} else {
		return w.Minion.Output(calldepth+1, info.Output(w.format))
	}
}

func colorString(color int) string {
	return fmt.Sprintf("\033[%dm", color)
}

func initColors() {
	colors = map[LogLevel]string{
		CriticalLevel: colorString(Magenta),
		ErrorLevel:    colorString(Red),
		WarningLevel:  colorString(Yellow),
		NoticeLevel:   colorString(Green),
		DebugLevel:    colorString(Cyan),
		InfoLevel:     colorString(White),
	}
}

func initFormatPlaceholders() {
	phfs = map[string]string{
		"%{id}":       "%[1]d",
		"%{time}":     "%[2]s",
		"%{module}":   "%[3]s",
		"%{filename}": "%[4]s",
		"%{file}":     "%[4]s",
		"%{line}":     "%[5]d",
		"%{level}":    "%[6]s",
		"%{lvl}":      "%.3[6]s",
		"%{message}":  "%[7]s",
	}
}

func New(module string, level LogLevel, color bool, out io.Writer) (Logger, error) {
	if module == "" {
		module = "DEFAULT"
	} else {
		module = strings.ToUpper(module)
	}
	if out == nil {
		out = os.Stdout
	}
	newWorker := NewWorker("", 0, color, out)
	newWorker.SetLogLevel(level)
	return &logger{Module: module, worker: newWorker}, nil
}

func (l *logger) Log(lvl LogLevel, message string) {
	l.logInternal(lvl, message, 2)
}

func (l *logger) logInternal(lvl LogLevel, message string, pos int) {
	now, _ := clock.Now(clock.Local)
	_, filename, line, _ := runtime.Caller(pos)
	filename = path.Base(filename)
	info := &Info{
		Id:       atomic.AddUint64(&logNo, 1),
		Time:     clock.Format(now, clock.FormatRFC3339),
		Module:   l.Module,
		Level:    lvl,
		Message:  message,
		Filename: filename,
		Line:     line,
	}
	l.worker.Log(lvl, 2, info)
}

func (l *logger) Fatal(message string) {
	l.logInternal(CriticalLevel, message, 2)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, a ...interface{}) {
	l.logInternal(CriticalLevel, fmt.Sprintf(format, a...), 2)
	os.Exit(1)
}

func (l *logger) Panic(message string) {
	l.logInternal(CriticalLevel, message, 2)
	panic(message)
}

func (l *logger) Panicf(format string, a ...interface{}) {
	l.logInternal(CriticalLevel, fmt.Sprintf(format, a...), 2)
	panic(fmt.Sprintf(format, a...))
}

func (l *logger) Critical(message string) {
	l.logInternal(CriticalLevel, message, 2)
}

func (l *logger) Criticalf(format string, a ...interface{}) {
	l.logInternal(CriticalLevel, fmt.Sprintf(format, a...), 2)
}

func (l *logger) Error(message string) {
	l.logInternal(ErrorLevel, message, 2)
}

func (l *logger) Errorf(format string, a ...interface{}) {
	l.logInternal(ErrorLevel, fmt.Sprintf(format, a...), 2)
}

func (l *logger) Warning(message string) {
	l.logInternal(WarningLevel, message, 2)
}

func (l *logger) Warningf(format string, a ...interface{}) {
	l.logInternal(WarningLevel, fmt.Sprintf(format, a...), 2)
}

func (l *logger) Notice(message string) {
	l.logInternal(NoticeLevel, message, 2)
}

func (l *logger) Noticef(format string, a ...interface{}) {
	l.logInternal(NoticeLevel, fmt.Sprintf(format, a...), 2)
}

func (l *logger) Info(message string) {
	l.logInternal(InfoLevel, message, 2)
}

func (l *logger) Infof(format string, a ...interface{}) {
	l.logInternal(InfoLevel, fmt.Sprintf(format, a...), 2)
}

func (l *logger) Debug(message string) {
	l.logInternal(DebugLevel, message, 2)
}

func (l *logger) Debugf(format string, a ...interface{}) {
	l.logInternal(DebugLevel, fmt.Sprintf(format, a...), 2)
}

func (l *logger) StackAsError(message string) {
	if message == "" {
		message = "Stack info"
	}
	message += "\n"
	l.logInternal(ErrorLevel, message+Stack(), 2)
}

func (l *logger) StackAsCritical(message string) {
	if message == "" {
		message = "Stack info"
	}
	message += "\n"
	l.logInternal(CriticalLevel, message+Stack(), 2)
}

func Stack() string {
	buf := make([]byte, 1000000)
	runtime.Stack(buf, false)
	return string(buf)
}

func (r *Info) logLevelString() string {
	logLevels := [...]string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"NOTICE",
		"INFO",
		"DEBUG",
	}
	return logLevels[r.Level-1]
}
