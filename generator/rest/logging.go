package rest

import "golang.org/x/net/context"

// Logger offers an environment agnostic interface for logging
type Logger interface {
	Debug(c context.Context, format string, args ...interface{})
	Info(c context.Context, format string, args ...interface{})
	Warning(c context.Context, format string, args ...interface{})
	Error(c context.Context, format string, args ...interface{})
	Critical(c context.Context, format string, args ...interface{})
}
