package main

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// Logger implements the default HTTP handler interface.
type Logger struct {
	h http.Handler
	w io.Writer
}

// NewLogger creates an instance of the HTTP request catcher.
func NewLogger(handler http.Handler, out io.Writer) http.Handler {
	return &Logger{h: handler, w: out}
}

// ServeHTTP calls f(w, r).
func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr

	if mark := strings.LastIndex(addr, ":"); mark != -1 {
		addr = addr[0:mark]
	}

	record := &AccessLog{
		ResponseWriter: rw,

		addr:     addr,
		time:     time.Time{},
		method:   r.Method,
		uri:      r.RequestURI,
		protocol: r.Proto,
		status:   http.StatusOK,
		elapsed:  time.Duration(0),
	}

	startTime := time.Now()
	l.h.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsed = finishTime.Sub(startTime)
	record.Log(l.w)
}
