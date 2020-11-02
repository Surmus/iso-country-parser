package internal

import (
	"github.com/sirupsen/logrus"
)

// CLIFormatter is logrus log entry formatter for outputting log messages into CLI
type CLIFormatter struct {
}

// Format transforms log entry suitable form for CLI output by omitting timestamp and context log entry components
func (f *CLIFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}
