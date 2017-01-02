package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LogRecord struct {
	http.ResponseWriter

	ip            string
	time          time.Time
	method        string
	uri           string
	protocol      string
	status        int
	responseBytes int64
	elapsedTime   time.Duration
}

func (r *LogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)

	fmt.Fprintf(out,
		FormatPattern,
		r.ip,
		timeFormatted,
		requestLine,
		r.status,
		r.responseBytes,
		r.elapsedTime.Seconds())
}

func (r *LogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type LoggingHandler struct {
	handler http.Handler
	out     io.Writer
}

func NewLoggingHandler(handler http.Handler, out io.Writer) http.Handler {
	return &LoggingHandler{handler, out}
}

func (h *LoggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &LogRecord{
		ResponseWriter: rw,
		ip:             clientIP,
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
