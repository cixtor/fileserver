package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LoggingHandler implements the default HTTP handler interface.
type LoggingHandler struct {
	handler http.Handler
	out     io.Writer
}

// LogRecord defines which data will be collected from the requests.
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

// NewLoggingHandler creates an instance of the HTTP request catcher.
func NewLoggingHandler(handler http.Handler, out io.Writer) http.Handler {
	return &LoggingHandler{handler, out}
}

// Log sends the access log string to the available writer.
func (r *LogRecord) Log(out io.Writer) {
	fmt.Fprintf(
		out,
		"%s - - [%s] \"%s %s %s\" %d %d %.4fs\n",
		r.ip,
		r.time.Format(`02/Jan/2006:15:04:05 -0700`),
		r.method,
		r.uri,
		r.protocol,
		r.status,
		r.responseBytes,
		r.elapsedTime.Seconds(),
	)
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *LogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// ServeHTTP calls f(w, r).
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
