// Package log provides a logging facility.
package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

type Logger interface {
	Emergency(format string, args ...interface{}) error
	Alert(format string, args ...interface{}) error
	Crit(format string, args ...interface{}) error
	Err(format string, args ...interface{}) error
	Warn(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

type level int

const (
	emergency level = iota
	alert
	crit
	err
	warn
	notice
	info
	debug
)

func (l level) String() string {
	return levels[l]
}

var levels = map[level]string{
	emergency: "emergency",
	alert:     "alert",
	crit:      "crit",
	err:       "err",
	warn:      "warn",
	notice:    "notice",
	info:      "info",
	debug:     "debug",
}

func NewLogger(levelStr string, timeFormat string, writer io.Writer) (Logger, error) {
	lvl, err := parseLevel(levelStr)
	if err != nil {
		return nil, err
	}
	return &logger{lvl, timeFormat, writer}, nil
}

func parseLevel(lvl string) (level, error) {
	for l, name := range levels {
		if name == lvl {
			return l, nil
		}
	}
	return 0, fmt.Errorf("unknown level: %s", lvl)
}

type logger struct {
	level      level
	timeFormat string
	writer     io.Writer
}

func (l *logger) Emergency(format string, args ...interface{}) error {
	return l.logError(emergency, format, args...)
}

func (l *logger) Alert(format string, args ...interface{}) error {
	return l.logError(alert, format, args...)
}

func (l *logger) Crit(format string, args ...interface{}) error {
	return l.logError(crit, format, args...)
}

func (l *logger) Err(format string, args ...interface{}) error {
	return l.logError(err, format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	l.log(warn, format, args...)
}

func (l *logger) Notice(format string, args ...interface{}) {
	l.log(notice, format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.log(info, format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log(debug, format, args...)
}

func (l *logger) logError(lvl level, format string, args ...interface{}) error {
	now := time.Now()
	if lvl <= l.level {
		l.fprintf(l.writer, now, lvl, format, args...)
	}
	var b bytes.Buffer
	l.fprintf(&b, now, lvl, format, args...)
	return errors.New(b.String())
}

func (l *logger) log(lvl level, format string, args ...interface{}) {
	if lvl > l.level {
		return
	}

	l.fprintf(l.writer, time.Now(), lvl, format, args...)
}

func (l *logger) fprintf(w io.Writer, t time.Time, lvl level, format string, args ...interface{}) {
	format = fmt.Sprintf("%s [%s] %s\n", t.Format(l.timeFormat), levels[lvl], format)
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		fmt.Printf("log error: %s: could not write to: %#v", err, w)
	}
}
