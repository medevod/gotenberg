package core

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

func ModifyLogger(level, format string) error {
	lvl, err := newLogLevel(level)
	if err != nil {
		return err
	}

	encoder, err := newLogEncoder(format)
	if err != nil {
		return err
	}

	loggerMu.RLock()
	defer loggerMu.RUnlock()

	core := zapcore.NewCore(encoder, os.Stderr, lvl)
	logger = zap.New(core)
	defer logger.Sync()

	return nil
}

func newLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "error":
		return zap.ErrorLevel, nil
	case "warn":
		return zap.WarnLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "debug":
		return zap.DebugLevel, nil
	default:
		return -2, fmt.Errorf("'%s' is not a recognized log level", level)
	}
}

func newLogEncoder(format string) (zapcore.Encoder, error) {
	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))
	encCfg := zap.NewProductionEncoderConfig()

	if isTerminal {
		// If interactive terminal, make output more human-readable by default.
		// Credits: https://github.com/caddyserver/caddy/blob/v2.1.1/logging.go#L671.
		encCfg.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(ts.UTC().Format("2006/01/02 15:04:05.000"))
		}

		if format == "text" || format == "auto" {
			encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	}

	if format == "auto" && isTerminal {
		format = "text"
	} else if format == "auto" {
		format = "json"
	}

	switch format {
	case "text":
		return zapcore.NewConsoleEncoder(encCfg), nil
	case "json":
		return zapcore.NewJSONEncoder(encCfg), nil
	default:
		return nil, fmt.Errorf("'%s' is not a recognized log format", format)
	}
}

func Log() *zap.Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()

	return logger
}

var (
	logger, _ = zap.NewProduction()
	loggerMu  sync.RWMutex
)
