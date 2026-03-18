package logger

import (
	"log"

	"github.com/zakiverse/zakiverse-api/core/cst"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(deployMode string) {
	config := zap.NewDevelopmentConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.MessageKey = "message"

	if deployMode == cst.DeployModeProduction {
		config = zap.NewProductionConfig()
		encoderConfig.StacktraceKey = ""
	}

	config.EncoderConfig = encoderConfig

	var err error
	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("[Error]->Failed to initialize logger : %s", err)
	}

	logger.Info("Initialized logger successfully")
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warning(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Field(key string, value any) zap.Field {
	return zap.Any(key, value)
}
