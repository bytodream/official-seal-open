// Package utils contains common properties used by the subpackages.
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// utcTimeLayout is the layout used by UtcTimeFormat.
const utcTimeLayout = "2 January 2006 at 15:04"

// UtcTimeFormat returns a formatted string describing a UTC timestamp.
func UtcTimeFormat(t time.Time) string { return t.Format(utcTimeLayout) + " UTC" }

// LogLevel is the verbosity level for logging.
type LogLevel int

// Log levels, ordered from the least to the most verbose.
const (
	LogLevelError   LogLevel = 0
	LogLevelWarning LogLevel = 1
	LogLevelNotice  LogLevel = 2
	LogLevelInfo    LogLevel = 3
	LogLevelDebug   LogLevel = 4
	LogLevelTrace   LogLevel = 5
)

// Logger is the logging interface used across the OS packages.
type Logger interface {
	Error(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	IsLogLevel(logLevel LogLevel) bool
}

// Log is the global logger. It defaults to a no-op logger.
var Log Logger = DummyLogger{}

// SetLogger sets 'logger' to be used by the OS library.
func SetLogger(logger Logger) { Log = logger }

// DummyLogger is a logger that does nothing.
type DummyLogger struct{}

// IsLogLevel always returns true, so callers do not skip log preparation.
func (DummyLogger) IsLogLevel(level LogLevel) bool { return true }

// Error does nothing for dummy logger.
func (DummyLogger) Error(format string, args ...interface{}) {}

// Warning does nothing for dummy logger.
func (DummyLogger) Warning(format string, args ...interface{}) {}

// Notice does nothing for dummy logger.
func (DummyLogger) Notice(format string, args ...interface{}) {}

// Info does nothing for dummy logger.
func (DummyLogger) Info(format string, args ...interface{}) {}

// Debug does nothing for dummy logger.
func (DummyLogger) Debug(format string, args ...interface{}) {}

// Trace does nothing for dummy logger.
func (DummyLogger) Trace(format string, args ...interface{}) {}

// ConsoleLogger is a logger that writes logs to the 'os.Stdout'.
type ConsoleLogger struct{ LogLevel LogLevel }

// NewConsoleLogger creates new console logger.
func NewConsoleLogger(logLevel LogLevel) *ConsoleLogger { return &ConsoleLogger{LogLevel: logLevel} }

// IsLogLevel returns true if log level is greater or equal than `level`.
// Can be used to avoid resource intensive calls to loggers.
func (consoleLogger ConsoleLogger) IsLogLevel(level LogLevel) bool {
	return consoleLogger.LogLevel >= level
}

// Error logs error message.
func (consoleLogger ConsoleLogger) Error(format string, args ...interface{}) {
	if consoleLogger.LogLevel >= LogLevelError {
		consoleLogger.output(os.Stdout, "[ERROR] ", format, args...)
	}
}

// Warning logs warning message.
func (consoleLogger ConsoleLogger) Warning(format string, args ...interface{}) {
	if consoleLogger.LogLevel >= LogLevelWarning {
		consoleLogger.output(os.Stdout, "[WARNING] ", format, args...)
	}
}

// Notice logs notice message.
func (consoleLogger ConsoleLogger) Notice(format string, args ...interface{}) {
	if consoleLogger.LogLevel >= LogLevelNotice {
		consoleLogger.output(os.Stdout, "[NOTICE] ", format, args...)
	}
}

// Info logs info message.
func (consoleLogger ConsoleLogger) Info(format string, args ...interface{}) {
	if consoleLogger.LogLevel >= LogLevelInfo {
		consoleLogger.output(os.Stdout, "[INFO] ", format, args...)
	}
}

// Debug logs debug message.
func (consoleLogger ConsoleLogger) Debug(format string, args ...interface{}) {
	if consoleLogger.LogLevel >= LogLevelDebug {
		consoleLogger.output(os.Stdout, "[DEBUG] ", format, args...)
	}
}

// Trace logs trace message.
func (consoleLogger ConsoleLogger) Trace(format string, args ...interface{}) {
	if consoleLogger.LogLevel >= LogLevelTrace {
		consoleLogger.output(os.Stdout, "[TRACE] ", format, args...)
	}
}

// output forwards a formatted message to writer.
// 保持两层调用深度，确保 message 中 runtime.Caller(3) 定位到用户代码。
func (consoleLogger ConsoleLogger) output(writer io.Writer, level string, format string, args ...interface{}) {
	message(writer, level, format, args...)
}

// WriterLogger is the logger that writes data to the Output writer.
type WriterLogger struct {
	LogLevel LogLevel
	Output   io.Writer
}

// NewWriterLogger creates new 'writer' logger.
func NewWriterLogger(logLevel LogLevel, writer io.Writer) *WriterLogger {
	writerLogger := WriterLogger{Output: writer, LogLevel: logLevel}
	return &writerLogger
}

// IsLogLevel returns true if log level is greater or equal than `level`.
// Can be used to avoid resource intensive calls to loggers.
func (writerLogger WriterLogger) IsLogLevel(level LogLevel) bool {
	return writerLogger.LogLevel >= level
}

// Error logs error message.
func (writerLogger WriterLogger) Error(format string, args ...interface{}) {
	if writerLogger.LogLevel >= LogLevelError {
		writerLogger.logToWriter(writerLogger.Output, "[ERROR] ", format, args...)
	}
}

// Warning logs warning message.
func (writerLogger WriterLogger) Warning(format string, args ...interface{}) {
	if writerLogger.LogLevel >= LogLevelWarning {
		writerLogger.logToWriter(writerLogger.Output, "[WARNING] ", format, args...)
	}
}

// Notice logs notice message.
func (writerLogger WriterLogger) Notice(format string, args ...interface{}) {
	if writerLogger.LogLevel >= LogLevelNotice {
		writerLogger.logToWriter(writerLogger.Output, "[NOTICE] ", format, args...)
	}
}

// Info logs info message.
func (writerLogger WriterLogger) Info(format string, args ...interface{}) {
	if writerLogger.LogLevel >= LogLevelInfo {
		writerLogger.logToWriter(writerLogger.Output, "[INFO] ", format, args...)
	}
}

// Debug logs debug message.
func (writerLogger WriterLogger) Debug(format string, args ...interface{}) {
	if writerLogger.LogLevel >= LogLevelDebug {
		writerLogger.logToWriter(writerLogger.Output, "[DEBUG] ", format, args...)
	}
}

// Trace logs trace message.
func (writerLogger WriterLogger) Trace(format string, args ...interface{}) {
	if writerLogger.LogLevel >= LogLevelTrace {
		writerLogger.logToWriter(writerLogger.Output, "[TRACE] ", format, args...)
	}
}

// logToWriter forwards a formatted message to writer.
// 保持两层调用深度，确保 message 中 runtime.Caller(3) 定位到用户代码。
func (writerLogger WriterLogger) logToWriter(writer io.Writer, level string, format string, args ...interface{}) {
	message(writer, level, format, args)
}

// message writes a log entry prefixed with the log level and the caller's file:line.
func message(writer io.Writer, level string, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}
	entry := fmt.Sprintf("%s %s:%d ", level, file, line) + format + "\n"
	fmt.Fprintf(writer, entry, args...)
}
