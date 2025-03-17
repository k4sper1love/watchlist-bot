// Package logger provides a thread-safe logging utility for managing user-specific loggers.
//
// It supports creating, managing, and writing logs to separate files for each user,
// ensuring efficient log rotation and customizable logging behavior.
//
// Features:
// - User-Specific Logging: Each user has a dedicated log file.
// - Thread Safety: All logging operations use mutexes for safe concurrent access.
// - Log Rotation: Uses the "lumberjack" library for automatic log file rotation based on size, backups, and age.
// - Customizable Prefix: Allows setting a custom prefix for log messages.
// - Flexible Logging Methods: Provides `Printf`, `Print`, and `Println` for structured logging.
//
// Usage:
//
// To log messages for a specific user, call `GetLogger(userID)` to retrieve or create a logger.
//
// Example:
//
//	logger := logger.GetLogger(userID)
//	logger.Printf("User %d performed an action", userID)
//
// Log Rotation:
//
// This package uses the "lumberjack" library for log rotation, based on:
// - Maximum file size (MB).
// - Maximum number of backup files.
// - Maximum age of log files (days).
// - Compression of old logs.
//
// Thread Safety:
//
// - All logging operations use `sync.Mutex` to prevent race conditions in concurrent writes.
// - Loggers are stored in a global map protected by `sync.RWMutex`.
//
// Logger Management:
//
// - Loggers are stored in a global map keyed by user ID.
// - If a logger does not exist for a user, it is created dynamically.
//
// This package simplifies logging in multi-user applications, ensuring organized and manageable log files.
package logger
