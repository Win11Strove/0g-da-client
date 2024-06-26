package logging

import (
	"fmt"
	"os"

	"github.com/0glabs/0g-da-client/common"
	"github.com/ethereum/go-ethereum/log"
)

type Logger struct {
	log.Logger
}

func (l *Logger) New(ctx ...interface{}) common.Logger {
	return &Logger{Logger: l.Logger.New(ctx...)}
}

func (l *Logger) SetHandler(h log.Handler) {
	l.Logger.SetHandler(h)
}

// GetLogger returns a logger with the specified configuration.
func GetLogger(cfg Config) (common.Logger, error) {
	fileLevel, err := log.LvlFromString(cfg.FileLevel)
	if err != nil {
		return nil, err
	}
	stdLevel, err := log.LvlFromString(cfg.StdLevel)
	if err != nil {
		return nil, err
	}

	logger := &Logger{Logger: log.New()}
	// This is required to print locations of log calls
	// This was recently added in this PR: https://github.com/ethereum/go-ethereum/pull/28069/files
	// where the default behavior was changed to not print origins
	// This was due to it being very expensive to compute origins
	// We should evaluate enabling/disabling this based on the flag
	log.PrintOrigins(true)
	stdh := log.StreamHandler(os.Stdout, log.TerminalFormat(false))
	stdHandler := log.CallerFileHandler(log.LvlFilterHandler(stdLevel, stdh))
	if cfg.Path != "" {
		fh, err := log.FileHandler(cfg.Path, log.LogfmtFormat())
		if err != nil {
			return nil, err
		}
		fileHandler := log.LvlFilterHandler(fileLevel, fh)
		logger.SetHandler(log.MultiHandler(fileHandler, stdHandler))
	} else {
		logger.SetHandler(stdHandler)
	}
	return logger, nil
}

func (l *Logger) Fatal(msg string, ctx ...interface{}) {
	l.Crit(msg, ctx...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Debug(fmt.Sprintf(template, args...))
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.Info(fmt.Sprintf(template, args...))
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Warn(fmt.Sprintf(template, args...))
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Error(fmt.Sprintf(template, args...))
}

func (l *Logger) Critf(template string, args ...interface{}) {
	l.Crit(fmt.Sprintf(template, args...))
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Crit(fmt.Sprintf(template, args...))
}
