// Package middleware provides HTTP middleware for authentication, authorization, and logging.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// LogEntry represents a single HTTP request log entry.
type LogEntry interface {
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra any)
	Panic(v any, stack []byte)
}

// LogFormatter creates new LogEntry instances for each request.
type LogFormatter interface {
	NewLogEntry(r *http.Request) LogEntry
}

// WrapResponseWriter extends http.ResponseWriter to capture status and bytes written.
type WrapResponseWriter interface {
	http.ResponseWriter
	Status() int
	BytesWritten() int
	Unwrap() http.ResponseWriter
}

// wrapResponseWriter implements WrapResponseWriter for tracking response status and size.
type wrapResponseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

// NewWrapResponseWriter creates a new WrapResponseWriter for the given ResponseWriter.
func NewWrapResponseWriter(w http.ResponseWriter, protoMajor int) WrapResponseWriter {
	return &wrapResponseWriter{ResponseWriter: w, status: 200}
}

// WriteHeader sets the HTTP status code for the response.
func (w *wrapResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

// Write writes the response body and tracks the number of bytes written.
func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// Status returns the HTTP status code written to the response.
func (w *wrapResponseWriter) Status() int { return w.status }

// BytesWritten returns the number of bytes written to the response.
func (w *wrapResponseWriter) BytesWritten() int { return w.bytes }

// Unwrap returns the underlying http.ResponseWriter.
func (w *wrapResponseWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// Flush implements http.Flusher for streaming and SSE support.
func (w *wrapResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// SlogLogEntry implements LogEntry for slog-based request logging.
type SlogLogEntry struct {
	start  time.Time
	method string
	url    string
	remote string
	reqID  string
	logger *slog.Logger
}

// slogLogFormatter implements LogFormatter for slog-based logging with a custom logger.
type slogLogFormatter struct {
	Logger *slog.Logger
}

// NewLogEntry creates a new SlogLogEntry for the given request, using the provided logger.
func (f *slogLogFormatter) NewLogEntry(r *http.Request) LogEntry {
	return &SlogLogEntry{
		start:  time.Now(),
		method: r.Method,
		url:    r.URL.String(),
		remote: r.RemoteAddr,
		reqID:  uuid.New().String(),
		logger: f.Logger,
	}
}

// Write logs the request/response details using slog.
func (e *SlogLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	}
	e.logger.Log(context.Background(), level, "request",
		"method", e.method,
		"url", e.url,
		"remote", e.remote,
		"status", status,
		"bytes", bytes,
		"duration_ms", elapsed.Milliseconds(),
		"request_id", e.reqID,
	)
}

// Panic logs a panic that occurred during request handling.
func (e *SlogLogEntry) Panic(v any, stack []byte) {
	e.logger.Error("panic in handler", "panic", v, "stack", string(stack), "request_id", e.reqID)
}

// RequestLogger returns middleware that logs HTTP requests using the provided LogFormatter.
// It recovers from panics and logs request/response details.
func RequestLogger(f LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := f.NewLogEntry(r)
			ww := NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				if recov := recover(); recov != nil {
					stack := make([]byte, 4096)
					n := runtime.Stack(stack, false)
					entry.Panic(recov, stack[:n])
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()
			next.ServeHTTP(ww, r)
		})
	}
}

// Slog returns middleware that logs HTTP requests using slog with the provided logger.
func Slog(logger *slog.Logger) func(http.Handler) http.Handler {
	return RequestLogger(&slogLogFormatter{Logger: logger})
}
