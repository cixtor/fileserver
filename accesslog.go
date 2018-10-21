package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// AccessLog defines which data will be collected from the requests.
type AccessLog struct {
	http.ResponseWriter

	addr      string
	time      time.Time
	method    string
	uri       string
	protocol  string
	status    int
	totalSize int64
	elapsed   time.Duration
}

// Log sends the access log string to the available writer.
func (r *AccessLog) Log(out io.Writer) {
	fmt.Fprintf(
		out,
		"%s - - [%s] \"%s %s %s\" %d %d %.4fs\n",
		r.addr,
		r.time.Format(`02/Jan/2006:15:04:05 -0700`),
		r.method,
		r.uri,
		r.protocol,
		r.status,
		r.totalSize,
		r.elapsed.Seconds(),
	)
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *AccessLog) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.totalSize += int64(written)
	return written, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (r *AccessLog) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
