package contextLogger

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxKey is used as the context key for the logger.
type ctxKey struct{}

// NewLogger creates a new zap.Logger based on the environment.
func NewLogger(env string) (*zap.Logger, error) {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.Development = false
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Development = true

		// Custom colorized level encoder
		cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			var levelColor string
			switch level {
			case zapcore.DebugLevel:
				levelColor = "\x1b[37mDEBUG\x1b[0m" // White
			case zapcore.InfoLevel:
				levelColor = "\x1b[32mINFO\x1b[0m" // Green
			case zapcore.WarnLevel:
				levelColor = "\x1b[33mWARN\x1b[0m" // Yellow
			case zapcore.ErrorLevel:
				levelColor = "\x1b[31mERROR\x1b[0m" // Red
			case zapcore.DPanicLevel:
				levelColor = "\x1b[35mDPANIC\x1b[0m" // Magenta
			case zapcore.PanicLevel:
				levelColor = "\x1b[41mPANIC\x1b[0m" // Red background
			case zapcore.FatalLevel:
				levelColor = "\x1b[41;97mFATAL\x1b[0m" // Red bg + white text
			default:
				levelColor = "\x1b[37mUNKNOWN\x1b[0m"
			}
			enc.AppendString(levelColor)
		}

		// Cyan-colored caller encoder
		cfg.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("\x1b[36m" + caller.TrimmedPath() + "\x1b[0m") // Cyan
		}
	}

	l, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Get retrieves the singleton zap.Logger instance.
func Get(env string) *zap.Logger {
	var loggerInstance *zap.Logger
	var err error
	var once sync.Once

	once.Do(func() {
		loggerInstance, err = NewLogger(env)
		if err != nil {
			loggerInstance = zap.NewNop()
		}
	})
	return loggerInstance
}

// FromCtx retrieves the *zap.Logger associated with the provided context.
// If no logger is found, it returns the global logger (or a no-op logger if neither exists).
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else {
		log.Error().Msg("No logger found in context, using nil logger")
		return zap.NewNop()
	}
}

// WithCtx associates the provided *zap.Logger with the given context.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}
