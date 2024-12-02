package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger инициализирует глобальный логгер
func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Цветные уровни
	config.EncoderConfig.TimeKey = "timestamp"                          // Время
	config.EncoderConfig.CallerKey = "caller"                           // Файл и строка
	config.EncoderConfig.MessageKey = "message"                         // Сообщение
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Формат времени

	baseLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Добавляем AddCallerSkip, чтобы исключить путь log.go из вызова
	logger = baseLogger.WithOptions(zap.AddCallerSkip(1)) // Пропускаем 1 уровень вызова
}

// GetLogger возвращает текущий экземпляр логгера
func GetLogger() *zap.Logger {
	if logger == nil {
		panic("Logger is not initialized. Call InitLogger() first.")
	}
	return logger
}

// SyncLogger синхронизирует буфер логгера (например, для flush)
func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// Debug обёртка для logger.Debug
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info обёртка для logger.Info
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn обёртка для logger.Warn
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error обёртка для logger.Error
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal обёртка для logger.Fatal
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
