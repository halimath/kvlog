package kvlog

import (
	"net/http"
	"time"
)

type accessLogHandler struct {
	logger   *Logger
	delegate http.Handler
}

type responseWriterWrapper struct {
	w          http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	return w.w.Write(data)
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (l *accessLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	wrapper := &responseWriterWrapper{
		w:          w,
		statusCode: 200,
	}

	l.delegate.ServeHTTP(wrapper, r)

	requestTime := time.Now().Sub(startTime)
	l.logger.Info(KV("event", "request"), KV("method", r.Method), KV("url", r.URL), KV("status", wrapper.statusCode), KV("duration", requestTime))
}

// Handler returns a http.Handler that acts as an access log middleware
func Handler(l *Logger, h http.Handler) http.Handler {
	return &accessLogHandler{
		logger:   l,
		delegate: h,
	}
}
