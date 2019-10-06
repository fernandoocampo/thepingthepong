package logging

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

// Options holds the configuration to setup the logger
type Options struct {
	LogLevel  string
	LogFields logrus.Fields
	LogFormat string
}

// Handle is our wrapper around logrus.Entry
type Handle struct {
	*logrus.Entry
	LevelLabel string
}

// NewLogger creates new logrus logger for a package
func NewLogger(o Options) (*Handle, error) {

	switch strings.ToLower(o.LogLevel) {
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	switch strings.ToLower(o.LogFormat) {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logger := &Handle{logrus.StandardLogger().WithFields(o.LogFields), strings.ToLower(o.LogLevel)}

	return logger, nil

}

// SetFormat sets the output format of the logs
// either text or json
func (lh *Handle) SetFormat(format string) error {
	if format != "" {
		switch strings.ToLower(format) {
		case "text":
			lh.Logger.Formatter = &logrus.TextFormatter{}
			return nil
		case "json":
			lh.Logger.Formatter = &logrus.JSONFormatter{}
			return nil
		default:
			return errors.New("Unknown format, expected text,json")
		}
	}
	return errors.New("Unknown format, expected text,json")
}

// SetLevel sets the log level: warn, info or debug
func (lh *Handle) SetLevel(level string) error {
	if level != "" {
		switch strings.ToLower(level) {
		case "warn":
			lh.Level = logrus.WarnLevel
			return nil
		case "info":
			lh.Level = logrus.InfoLevel
			return nil
		case "debug":
			lh.Level = logrus.DebugLevel
			return nil
		default:
			return errors.New("Unknown LogLevel, expected Debug,Info or Warn")
		}
	}
	return errors.New("Unknown LogLevel, expected Debug,Info or Warn")
}
