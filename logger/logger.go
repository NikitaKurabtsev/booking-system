package logger

type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}
