package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger представляет интерфейс для логгера
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Close()
}

// ZapLogger реализует интерфейс Logger с помощью zap
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger создает новый экземпляр ZapLogger
func NewZapLogger(level string) (*ZapLogger, error) {
	logLevel := zap.NewAtomicLevel()

	// Установка уровня логирования
	switch level {
	case "debug":
		logLevel.SetLevel(zap.DebugLevel)
	case "info":
		logLevel.SetLevel(zap.InfoLevel)
	case "warn":
		logLevel.SetLevel(zap.WarnLevel)
	case "error":
		logLevel.SetLevel(zap.ErrorLevel)
	default:
		logLevel.SetLevel(zap.InfoLevel) // По умолчанию INFO
	}

	// Настройка кодировщика
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Создание консольного кодировщика
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Создание ядра логгера
	core := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout), // Вывод в stdout
		logLevel,
	)

	// Создание логгера
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return &ZapLogger{logger: zapLogger}, nil
}

// Info записывает информационное сообщение
func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Debug записывает отладочное сообщение
func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Warn записывает предупреждение
func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error записывает ошибку
func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal записывает критическую ошибку и завершает программу
func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Close закрывает логгер
func (l *ZapLogger) Close() {
	l.logger.Sync()
}
