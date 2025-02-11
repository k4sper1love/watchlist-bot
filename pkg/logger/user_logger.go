package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"sync"
)

type Wrapper struct {
	mu     sync.Mutex
	logger *log.Logger
}

var (
	loggers = make(map[int]*Wrapper)
	mutex   sync.Mutex
)

type Logger log.Logger

func GetLogger(userID int) *Wrapper {
	mutex.Lock()
	defer mutex.Unlock()

	if logger, exists := loggers[userID]; exists {
		return logger
	}

	logFile := fmt.Sprintf("logs/user_%d.log", userID)
	logger := log.New(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}, "", log.Ldate|log.Ltime|log.Lmsgprefix)

	loggerWrapper := &Wrapper{logger: logger}
	loggers[userID] = loggerWrapper

	return loggerWrapper
}

func (w *Wrapper) SetPrefix(prefix string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.SetPrefix(prefix)
}

func (w *Wrapper) Printf(format string, v ...interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Printf(format, v...)
}
