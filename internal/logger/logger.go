package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

type customEncoder struct {
	zapcore.Encoder
}

func (c *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	formatted := entry.Time.Format("2006/01/02 - 15:04:05") + " | " +
		entry.Level.CapitalString() + " | " +
		entry.Caller.TrimmedPath() + " | " +
		entry.Message

	buf, err := c.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	buf.Reset()
	buf.AppendString(formatted)
	return buf, nil
}

func newCustomEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &customEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		newCustomEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zap.DebugLevel,
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

// Info logs a message at InfoLevel. The message includes any fields passed.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof formats and logs a message at InfoLevel.
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf formats and logs a message at ErrorLevel.
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf formats and logs a message at DebugLevel.
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf formats and logs a message at WarnLevel.
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf formats and logs a message at PanicLevel.
func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}
