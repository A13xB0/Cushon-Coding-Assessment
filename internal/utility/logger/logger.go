package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SInitialises a new zap logger
func New(service, level string) (*zap.SugaredLogger, error) {
	config := ecszap.ECSCompatibleEncoderConfig(zap.NewProductionEncoderConfig())
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	defaultLogLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(defaultLogLevel)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel)
	logger := zap.New(ecszap.WrapCore(core), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	sugaredLogger := logger.Named(service).WithOptions(zap.AddCallerSkip(1)).Sugar()
	return sugaredLogger, nil
}
