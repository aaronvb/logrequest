package logrequest

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// LogRequest struct
// Pass values from middleware to this struct
type LogRequest struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Handler http.Handler
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

// ToLogger will print the Started and Completed request info to the passed logger
func (lr LogRequest) ToLogger(logger *log.Logger) {
	sw, completedDuration := lr.parseRequest()
	logger.Printf(`Started %s "%s" %s %s`, lr.Request.Method, lr.Request.URL.RequestURI(), lr.Request.RemoteAddr, lr.Request.Proto)
	logger.Printf("Completed %d in %s", sw.statusCode, completedDuration)
}

// ToString will return a map with the key 'started' and 'completed' that contain
// a string output for eacch
func (lr LogRequest) ToString() map[string]string {
	sw, completedDuration := lr.parseRequest()
	ts := make(map[string]string)
	ts["started"] = fmt.Sprintf(`Started %s "%s" %s %s`, lr.Request.Method, lr.Request.URL.RequestURI(), lr.Request.RemoteAddr, lr.Request.Proto)
	ts["completed"] = fmt.Sprintf("Completed %d in %s", sw.statusCode, completedDuration)

	return ts
}

// parseRequest will time the request and retrieve the status from the
// ResponseWriter. Returns the statusWriter struct and the duration
// of the request.
func (lr LogRequest) parseRequest() (statusWriter, time.Duration) {
	startTime := time.Now()
	sw := statusWriter{ResponseWriter: lr.Writer}
	lr.Handler.ServeHTTP(&sw, lr.Request)
	return sw, time.Now().Sub(startTime)
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = 200
	}

	n, err := w.ResponseWriter.Write(b)
	return n, err
}
