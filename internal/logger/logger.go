package logger

import (
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func Init(debug bool, lDir, lFile string) {
	logDir := lDir
	logFile := filepath.Join(logDir, lFile)

	// Ensure logs directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}

	// Lumberjack handles file rotation automatically
	rotator := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // MB per file
		MaxBackups: 5,    // keep last 5 files
		MaxAge:     30,   // days
		Compress:   true, // compress old logs
	}

	writer := zapcore.AddSync(io.MultiWriter(os.Stdout, rotator))

	var encoderCfg zapcore.EncoderConfig
	if debug {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if debug {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(encoder, writer, level)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
