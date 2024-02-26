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
	builder.Logger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	builder.Config(&logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
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

	var (
		level slog.Level
		msg   string
		attrs = make([]slog.Attr, 0, 2)
	)

	elapsed := time.Since(begin)
	attrs = append(attrs, slog.Duration("elapsed", elapsed))

	sql, rows := fc()
	attrs = append(attrs, slog.String("sql", sql))
	if rows != -1 {
		attrs = append(attrs, slog.Int64("rows", rows))
	}

	switch {
	case l.Config.LogLevel >= logger.Error && err != nil && (!errors.Is(err, logger.ErrRecordNotFound) || !l.Config.IgnoreRecordNotFoundError):
		level = slog.LevelError
		msg = "gorm: error"
		attrs = append(attrs, slog.String("error", err.Error()))
	case l.Config.LogLevel >= logger.Warn && l.Config.SlowThreshold != 0 && elapsed > l.Config.SlowThreshold:
		level = slog.LevelWarn
		msg = fmt.Sprintf("SLOW SQL >= %v", l.Config.SlowThreshold)
	case l.Config.LogLevel >= logger.Info:
		level = slog.LevelInfo
		msg = "gorm: info"
	default:
		return
	}

	l.Logger.LogAttrs(ctx, level, msg, attrs...)
}
