package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
)

// LogContext is a context key.
var LogContext = contextKey("request_logs")

// A LogCollector is used to collect logs.
type LogCollector struct {
	Messages []string
	lock     sync.Mutex
}

// Append appends msg into log context.
func (lc *LogCollector) Append(msg string) {
	lc.lock.Lock()
	lc.Messages = append(lc.Messages, msg)
	lc.lock.Unlock()
}

// Flush flushes collected logs.
func (lc *LogCollector) Flush() string {
	var buffer bytes.Buffer

	start := true
	for _, message := range lc.takeAll() {
		if start {
			start = false
		} else {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(message)
	}

	return buffer.String()
}

func (lc *LogCollector) takeAll() []string {
	lc.lock.Lock()
	messages := lc.Messages
	lc.Messages = nil
	lc.lock.Unlock()

	return messages
}

// Error logs the given v along with r in error log.
func Error(r *http.Request, v ...interface{}) {
	logx.WithContext(r.Context()).Error(v)
}

// Errorf logs the given v with format along with r in error log.
func Errorf(r *http.Request, format string, v ...interface{}) {
	logx.WithContext(r.Context()).Errorf(format, v)
}

// Info logs the given v along with r in access log.
func Info(r *http.Request, v ...interface{}) {
	appendLog(r, format(r, v...))
}

// Infof logs the given v with format along with r in access log.
func Infof(r *http.Request, format string, v ...interface{}) {
	appendLog(r, formatf(r, format, v...))
}

func appendLog(r *http.Request, message string) {
	logs := r.Context().Value(LogContext)
	if logs != nil {
		logs.(*LogCollector).Append(message)
	}
}

func format(r *http.Request, v ...interface{}) string {
	return formatWithReq(r, fmt.Sprint(v...))
}

func formatf(r *http.Request, format string, v ...interface{}) string {
	return formatWithReq(r, fmt.Sprintf(format, v...))
}

func formatWithReq(r *http.Request, v string) string {
	return fmt.Sprintf("(%s - %s) %s", r.RequestURI, httpx.GetRemoteAddr(r), v)
}

type contextKey string

func (c contextKey) String() string {
	return "rest/internal context key " + string(c)
}
