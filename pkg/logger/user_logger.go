package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"sync"
)

type Wrapper struct {
	logger *log.Logger
	mu     sync.Mutex
}

var (
	loggers = make(map[int]*Wrapper)
	mu      sync.RWMutex
)

func GetLogger(userID int) *Wrapper {
	mu.RLock()
	logger, exists := loggers[userID]
	mu.RUnlock()

	if exists {
		return logger
	}

	mu.Lock()
	defer mu.Unlock()

	if logger, exists = loggers[userID]; exists {
		return logger
	}

	logFile := fmt.Sprintf("logs/user_%d.log", userID)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	logger = &Wrapper{
		logger: log.New(lumberjackLogger, "", log.Ldate|log.Ltime|log.Lmsgprefix),
	}

	loggers[userID] = logger
	return logger
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

func (w *Wrapper) Print(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Print(msg)
}

func (w *Wrapper) Println(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Println(msg)
}
