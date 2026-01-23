package logx

import (
	"log"
	"strings"
)

type level int

const (
	Debug level = iota
	Info
	Warn
	Error
)

func ParseLevel(s string) level {
	switch strings.ToLower(s) {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		return Info
	}
}

type Logger struct {
	prefix string
	level  level
}

// New 创建Logger对象
func New(prefix string, level level) Logger {
	return Logger{prefix: prefix, level: level}
}

func (l Logger) Debug(msg string, kv ...any) { l.print(Debug, "debug", msg, kv...) }
func (l Logger) Info(msg string, kv ...any)  { l.print(Info, "info", msg, kv...) }
func (l Logger) Warn(msg string, kv ...any)  { l.print(Warn, "warn", msg, kv...) }
func (l Logger) Error(msg string, kv ...any) { l.print(Error, "error", msg, kv...) }

func (l Logger) print(lv level, lvStr string, msg string, kv ...any) {
	if lv < l.level {
		return
	}
	log.Println(append([]any{"level=" + lvStr, "msg=" + msg, l.prefix}, kv...)...)
}
