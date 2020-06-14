package logger

import (
	"errors"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger interface describes logging abstraction layer.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithPrefix(pfx string) Logger
}

// NOP provide no-ops implementation for Logger interface.
type NOP struct{}

func (nl NOP) Debug(args ...interface{})                 {}
func (nl NOP) Debugf(format string, args ...interface{}) {}
func (nl NOP) Info(args ...interface{})                  {}
func (nl NOP) Infof(format string, args ...interface{})  {}
func (nl NOP) Warn(args ...interface{})                  {}
func (nl NOP) Warnf(format string, args ...interface{})  {}
func (nl NOP) Error(args ...interface{})                 {}
func (nl NOP) Errorf(format string, args ...interface{}) {}
func (nl NOP) Fatal(args ...interface{})                 {}
func (nl NOP) Fatalf(format string, args ...interface{}) {}
func (nl NOP) WithPrefix(name string) Logger             { return NOP{} }

// LogMode defines different logging strategies.
type LogMode string

const (
	DevelopMode    LogMode = "develop"
	ProductionMode LogMode = "production"
)

// LogLevel represent different depths for logged informations about the system.
type LogLevel string

const (
	// Debug level logging, should be disabled in production mode.
	Debug LogLevel = "debug"
	// Info level logging, default logging setting
	Info LogLevel = "info"
	// Warn level logging used to mark possible issues in code execution.
	// For convenience "warn" spelling is also accepted.
	Warn LogLevel = "warning"
	// Error level logging logs errors
	Error LogLevel = "error"
	// Fatal level logging, log is recorded and application terminates after.
	Fatal LogLevel = "fatal"

	warnShort LogLevel = "warn"
)

func logLvlToZapLevel(ll LogLevel) zapcore.Level {
	switch ll {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn, warnShort:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Config represent high-level logger configuration.
type Config struct {
	Mode    LogMode
	Console ConsoleConfig
}

// ConsoleConfig groups configuration bits
type ConsoleConfig struct {
	Enabled bool
	Level   LogLevel
	Format  string
}

// DefaultConfig provides ready to use configuration for development mode.
// It is configured to support only console output with Debug logging
// level.
func DefaultConfig() Config {
	return Config{
		Mode: DevelopMode,
		Console: ConsoleConfig{
			Enabled: true,
			Level:   "debug",
			Format:  "console",
		},
	}
}

var emptyConfig Config

// New returns an instance of structured logger.
func New(cfg Config) Logger {
	if cfg == emptyConfig {
		return NOP{}
	}
	return new(cfg)
}

func new(cfg Config) Logger {

	switch cfg.Mode {
	case DevelopMode:
		lgr, err := developSugarZap(cfg)
		if err != nil {
			log.Fatal(err)
		}
		return wrapSugarZap(lgr)
	case ProductionMode:
		fallthrough
	default:
		log.Fatal("production logging mode is not implemented")
	}
	return NOP{}
}

type prefixableLogger struct {
	*zap.SugaredLogger
}

func (p *prefixableLogger) WithPrefix(pfx string) Logger {
	newLgr := prefixableLogger{
		SugaredLogger: p.Named(pfx),
	}
	return &newLgr
}

func wrapSugarZap(l *zap.SugaredLogger) *prefixableLogger {
	return &prefixableLogger{
		SugaredLogger: l,
	}
}

func developSugarZap(cfg Config) (*zap.SugaredLogger, error) {
	var cores []zapcore.Core
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Console.Enabled {
		var enc zapcore.Encoder
		if cfg.Console.Format == "json" {
			enc = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
			enc = zapcore.NewConsoleEncoder(encoderCfg)
		}
		ws := zapcore.Lock(os.Stdout)
		c := zapcore.NewCore(enc, ws, logLvlToZapLevel(cfg.Console.Level))
		cores = append(cores, c)
	}

	if len(cores) == 0 {
		return nil, errors.New("no logger")
	}

	// use Tee pattern to log to multiple sources.
	cc := zapcore.NewTee(cores...)
	lgr := zap.New(
		cc,
		zap.AddCaller(),
	)

	return lgr.Sugar(), nil
}
