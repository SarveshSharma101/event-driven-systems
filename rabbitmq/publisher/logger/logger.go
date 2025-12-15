package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Custom logger struct
type Logger struct {
	z     *zap.Logger
	sugar *zap.SugaredLogger
}

var (
	logInstance *Logger
	once        sync.Once
)

// InitLogger initializes the global logger instance
func InitLogger() {
	once.Do(func() {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		config := zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderConfig,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stdout"},
			InitialFields: map[string]any{
				"pid": os.Getpid(),
			},
		}

		z := zap.Must(config.Build())
		logInstance = &Logger{
			z:     z,
			sugar: z.Sugar(),
		}
	})
}

// Get returns the global logger instance
func Get() *Logger {
	if logInstance == nil {
		InitLogger()
	}
	return logInstance
}

//
// ------------- Logging methods -------------
//  Sugared logging (you don’t need zap.String()).
//

func (l *Logger) Info(msg string, fields ...any) {
	l.sugar.Infow(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...any) {
	l.sugar.Warnw(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...any) {
	l.sugar.Errorw(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...any) {
	l.sugar.Debugw(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...any) {
	l.sugar.Fatalw(msg, fields...)
}

//
// Structured logging version (optional)
//

func (l *Logger) InfoS(msg string, fields ...zap.Field) {
	l.z.Info(msg, fields...)
}

func (l *Logger) ErrorS(msg string, fields ...zap.Field) {
	l.z.Error(msg, fields...)
}

//
// Cleanup
//

func (l *Logger) Sync() {
	_ = l.z.Sync()
}
