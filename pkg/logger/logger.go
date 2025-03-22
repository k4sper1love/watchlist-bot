package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"os"
	"sync"
)

// Wrapper represents a thread-safe logger for a specific user.
// It uses a mutex to ensure safe concurrent access to the underlying logger.
type Wrapper struct {
	logger *log.Logger // The underlying logger instance.
	mu     sync.Mutex  // Mutex to ensure thread-safe operations.
}

// Global variables for managing loggers.
var (
	loggers = make(map[int]*Wrapper) // Map to store loggers keyed by user ID.
	mu      sync.RWMutex             // Read-write mutex for safe concurrent access to the loggers map.
)

// Get retrieves or creates a logger for the specified user ID.
// If a logger already exists for the user, it is returned. Otherwise, a new logger is created,
// configured with log rotation using the `lumberjack.Logger`, and stored in the global map.
func Get(userID int) *Wrapper {
	// Attempt to retrieve the logger from the map using a read lock.
	mu.RLock()
	logger, exists := loggers[userID]
	mu.RUnlock()

	if exists {
		return logger
	}

	// Acquire a write lock to safely create and store a new logger if it doesn't exist.
	mu.Lock()
	defer mu.Unlock()

	// Double-check if the logger was created while waiting for the lock.
	if logger, exists = loggers[userID]; exists {
		return logger
	}

	// Configure the log file path and log rotation settings.
	logFile := fmt.Sprintf("%s/users/user_%d.log", getLogDirectory(), userID)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile, // Path to the log file.
		MaxSize:    5,       // Maximum size of the log file in MB before rotation.
		MaxBackups: 3,       // Maximum number of old log files to retain.
		MaxAge:     7,       // Maximum age of old log files in days.
		Compress:   true,    // Compress old log files during rotation.
	}

	// Create a new logger with the specified settings.
	logger = &Wrapper{
		logger: log.New(lumberjackLogger, "", log.Ldate|log.Ltime|log.Lmsgprefix),
	}

	// Store the logger in the global map and return it.
	loggers[userID] = logger
	return logger
}

// SetPrefix sets a custom prefix for the logger's messages.
// The operation is thread-safe, protected by a mutex.
func (w *Wrapper) SetPrefix(prefix string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.SetPrefix(prefix)
}

// Printf formats and writes a log message using the specified format and arguments.
// The operation is thread-safe, protected by a mutex.
func (w *Wrapper) Printf(format string, v ...interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Printf(format, v...)
}

// Print writes a log message without formatting.
// The operation is thread-safe, protected by a mutex.
func (w *Wrapper) Print(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Print(msg)
}

// Println writes a log message followed by a newline character.
// The operation is thread-safe, protected by a mutex.
func (w *Wrapper) Println(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Println(msg)
}

// GetFilePath constructs the log file path for a user based on their ID.
// It checks if the log file exists and returns an error if it does not.
func GetFilePath(userID int) (string, error) {
	logFile := fmt.Sprintf("%s/users/user_%d.log", getLogDirectory(), userID)
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "", err
	}
	return logFile, nil
}

// getLogDirectory returns the log directory path.
// If the LOGS_DIR environment variable is set, it uses that value.
// Otherwise, it defaults to "./logs".
func getLogDirectory() string {
	if logsDir := os.Getenv("LOGS_DIR"); logsDir != "" {
		return logsDir
	}
	return "./logs" // Default log directory.
}
