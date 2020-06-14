package logger

import "fmt"

// EchoMiddlewareLogger is Echo Middleware compatible adapter.
type EchoMiddlewareLogger struct {
	DebugOnly bool
	lgr       Logger
}

// NewEchoMiddlewareLogger serve as adaptation layer to convert application
// logger to be compatible with echo framework middleware logging logic.
func NewEchoMiddlewareLogger(lgr Logger) *EchoMiddlewareLogger {
	return &EchoMiddlewareLogger{
		lgr: lgr,
	}
}

func (h *EchoMiddlewareLogger) Write(data []byte) (int, error) {
	if h.DebugOnly {
		h.lgr.Debug(fmt.Sprintf("%s", data))
		return 0, nil
	}
	h.lgr.Info(fmt.Sprintf("%s", data))
	return 0, nil
}
