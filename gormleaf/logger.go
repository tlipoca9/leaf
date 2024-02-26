package gormleaf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

type SlogLoggerBuilder struct {
	data SlogLogger
}

func NewSlogLoggerBuilder() *SlogLoggerBuilder {
	builder := &SlogLoggerBuilder{}
	builder.data.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return builder
}

func (b *SlogLoggerBuilder) Logger(logger *slog.Logger) *SlogLoggerBuilder {
	if !logger.Enabled(context.TODO(), slog.LevelDebug) {
		panic("slog logger must be enabled at debug level")
	}
	b.data.Logger = logger
	return b
}

func (b *SlogLoggerBuilder) Config(config *logger.Config) *SlogLoggerBuilder {
	b.data.Config = config
	return b
}

func (b *SlogLoggerBuilder) Build() logger.Interface {
	return &b.data
}

var _ logger.Interface = (*SlogLogger)(nil)

type SlogLogger struct {
	Logger *slog.Logger
	Config *logger.Config
}

// LogMode implements logger.Interface.
func (l *SlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.Config.LogLevel = level
	return l
}

// Info implements logger.Interface.
func (l *SlogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Info {
		l.Logger.InfoContext(ctx, fmt.Sprintf(msg, data...))
	}
}

// Warn implements logger.Interface.
func (l *SlogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Warn {
		l.Logger.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

// Error implements logger.Interface.
func (l *SlogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Error {
		l.Logger.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

// Trace implements logger.Interface.
func (l *SlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.Config.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Config.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.Config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Logger.ErrorContext(
				ctx, "",
				"error", err,
				"elapsed", elapsed,
				"sql", sql,
			)
		} else {
			l.Logger.ErrorContext(
				ctx, "",
				"error", err,
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
		}
	case l.Config.SlowThreshold != 0 && elapsed > l.Config.SlowThreshold && l.Config.LogLevel >= logger.Warn:
		sql, rows := fc()
		msg := fmt.Sprintf("SLOW SQL >= %v", l.Config.SlowThreshold)
		if rows == -1 {
			l.Logger.WarnContext(
				ctx, msg,
				"elapsed", elapsed,
				"sql", sql,
			)
		} else {
			l.Logger.WarnContext(
				ctx, msg,
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
		}
	case l.Config.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Logger.InfoContext(
				ctx, "",
				"elapsed", elapsed,
				"sql", sql,
			)
		} else {
			l.Logger.InfoContext(
				ctx, "",
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
		}
	}
}
