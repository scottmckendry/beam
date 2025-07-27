package middleware

import (
	"bytes"
	"net/http"
	"testing"
)

// mockResponseWriter is a minimal http.ResponseWriter for testing.
type mockResponseWriter struct {
	headers http.Header
	status  int
	buf     bytes.Buffer
}

func (m *mockResponseWriter) Header() http.Header         { return m.headers }
func (m *mockResponseWriter) Write(b []byte) (int, error) { return m.buf.Write(b) }
func (m *mockResponseWriter) WriteHeader(statusCode int)  { m.status = statusCode }

func TestWrapResponseWriter_StatusAndBytes(t *testing.T) {
	w := &mockResponseWriter{headers: make(http.Header)}
	ww := NewWrapResponseWriter(w, 1)

	if ww.Status() != 200 {
		t.Errorf("default status = %d, want 200", ww.Status())
	}

	ww.WriteHeader(404)
	if ww.Status() != 404 {
		t.Errorf("status after WriteHeader = %d, want 404", ww.Status())
	}

	ww.Write([]byte("hello"))
	if ww.BytesWritten() != 5 {
		t.Errorf("bytes written = %d, want 5", ww.BytesWritten())
	}
}

func TestWrapResponseWriter_Unwrap(t *testing.T) {
	w := &mockResponseWriter{headers: make(http.Header)}
	ww := NewWrapResponseWriter(w, 1)
	if ww.Unwrap() != w {
		t.Error("Unwrap() did not return original ResponseWriter")
	}
}

// flusher implements http.Flusher for testing Flush passthrough.
type flusher struct {
	*mockResponseWriter
	flushed bool
}

func (f *flusher) Flush() { f.flushed = true }

func TestWrapResponseWriter_Flush(t *testing.T) {
	f := &flusher{mockResponseWriter: &mockResponseWriter{headers: make(http.Header)}}
	ww := NewWrapResponseWriter(f, 1)
	if wwf, ok := ww.(interface{ Flush() }); ok {
		wwf.Flush()
		if !f.flushed {
			t.Error("Flush() did not call underlying Flush")
		}
	}
}
