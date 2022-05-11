package internalhttp

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	BytesLength int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.BytesLength += n

	return n, err
}

func loggingMiddleware(next http.Handler, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{w, 0, 0}
		next.ServeHTTP(writer, r)
		log.LogRequest(r, writer.StatusCode, writer.BytesLength)
	})
}
