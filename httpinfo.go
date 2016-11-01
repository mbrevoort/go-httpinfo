package httpinfo

import (
	"net/http"
	"time"
)

// HTTPInfo is wrapper of http.ResponseWriter that keeps track of its HTTP status
// code and body size
type httpInfo struct {
	w       http.ResponseWriter
	h       http.Handler
	status  int
	size    int
	reqSize int
	elapsed time.Duration
}

// HTTPInfo interface
type HTTPInfo interface {
	Header() http.Header
	Status() int
	Size() int
	ReqSize() int
	Elapsed() time.Duration
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// New returns a new HTTPInfo
//
// Simple http.Handler middleware that records HTTP status code, response size, and duration and makes
// the data available after `ServeHTTP` is finished. Requires no 3rd party dependencies.
//
// For example:
//
//    func FooMiddleware(h http.Handler) http.Handler {
//      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//        // wrap the request and track the response status code, size and duration
//        info := httpinfo.New(h)
//        info.ServeHTTP(w, r)
//
//        // log response
//        fmt.Printf("Request: %s %s %s %d %d (%d) ", r.Method, r.RequestURI, r.Proto, info.Status(), info.Size(), info.Elapsed())
//      })
//    }
func New(h http.Handler) HTTPInfo {
	l := httpInfo{h: h}
	return &l
}

// Header returns HTTP response headers
func (l *httpInfo) Header() http.Header {
	return l.w.Header()
}

func (l *httpInfo) Write(b []byte) (int, error) {
	if l.status == 0 {
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

// WriteHeader writes and records HTTP status code
func (l *httpInfo) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

// Status returns HTTP status code
func (l *httpInfo) Status() int {
	return l.status
}

// Size returns size of response
func (l *httpInfo) Size() int {
	return l.size
}

func (l *httpInfo) ReqSize() int {
	return l.reqSize
}

// Elapsed returns duration of time for the request
func (l *httpInfo) Elapsed() time.Duration {
	return l.elapsed
}

func (l *httpInfo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqLenChan := make(chan int)
	go computeApproximateRequestSize(req, reqLenChan)
	l.w = w
	t := time.Now()
	l.h.ServeHTTP(l, req)
	l.elapsed = time.Since(t)
	l.reqSize = <-reqLenChan
}

// Derived from https://github.com/prometheus/client_golang
func computeApproximateRequestSize(r *http.Request, out chan int) {
	s := len(r.URL.String())
	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	out <- s
}
