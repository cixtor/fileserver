package main

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// Logger implements the default HTTP handler interface.
type Logger struct {
	handler http.Handler
	out     io.Writer
}

// NewLogger creates an instance of the HTTP request catcher.
func NewLogger(handler http.Handler, out io.Writer) http.Handler {
	return &Logger{handler, out}
}

// ServeHTTP calls f(w, r).
func (h *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &AccessLog{
		ResponseWriter: rw,
		addr:           clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.out)
}
