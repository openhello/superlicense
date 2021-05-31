package log

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Glog *zap.Logger
	// change level use Level.SetLevel(zapcore.InfoLevel)
	Level = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func init() { // init for test
	Glog, _ = zap.NewProduction()
}

// ZapConfig for zap rotation
type ZapConfig struct {
	Filename   string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Level      string
	WithStdout bool
}

// InitZap config zap
func InitZap(c *ZapConfig) {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.

	os.MkdirAll(filepath.Dir(c.Filename), 0755)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
	})

	if c.Level == "debug" {
		Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		Level,
	)

	if c.WithStdout {
		Glog = zap.New(zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.Lock(os.Stdout),
				Level,
			),
			core,
		)).WithOptions(zap.AddCaller())
	} else {
		Glog = zap.New(core).WithOptions(zap.AddCaller())
	}
}
