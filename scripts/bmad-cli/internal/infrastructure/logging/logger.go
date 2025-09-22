package logging

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, err error, args ...interface{})
	Debug(msg string, args ...interface{})
}
