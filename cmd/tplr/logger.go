package main

import (
	"github.com/sebasmannem/tplr/internal"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.SugaredLogger
	atom zap.AtomicLevel
)

func initLogger() {

	atom = zap.NewAtomicLevel()
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= atom.Level()
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	//encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	log = zap.New(core).Sugar()
}

func initRemoteLoggers() {
	internal.InitLogger(log)
}

func enableDebug(debug bool) {
	if debug {
		atom.SetLevel(zap.DebugLevel)
	}
	log.Debug("Debug logging enabled by config")
}
