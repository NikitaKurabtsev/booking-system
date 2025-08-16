package logger

import (
	"log/slog"
	"os"
)

type SLogLogger struct {
	l *slog.Logger
}

func (s *SLogLogger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s *SLogLogger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}

func NewSLogger() Logger {
	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(loggerHandler)

	return &SLogLogger{logger}
}
