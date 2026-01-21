package logx

import "log"

type Logger struct {
	prefix string
}

// New 创建Logger对象
func New(prefix string) Logger {
	return Logger{prefix: prefix}
}

func (l Logger) Info(msg string, v ...any) {
	log.Println(append([]any{"level=info", "msg=" + msg, l.prefix}, v...)...)
}

func (l Logger) Error(msg string, v ...any) {
	log.Println(append([]any{"level=error", "msg=" + msg, l.prefix}, v...)...)
}
