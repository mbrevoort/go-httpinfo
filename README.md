Go HTTPInfo Middleware
----------------------

Simple http.Handler middleware that records HTTP status code, response size, and duration and makes
the data available after `ServeHTTP` is finished. Requires no 3rd party dependencies.

For example:

    func FooMiddleware(h http.Handler) http.Handler {
    	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    		// wrap the request and track the response status code, size and duration
    		info := httpinfo.New(h)
    		info.ServeHTTP(w, r)

    		// log response
        fmt.Printf("Request: %s %s %s %d %d (%d) ", r.Method, r.RequestURI, r.Proto, info.Status(), info.Size(), info.Elapsed())
    	})
    }


## Install

    go get github.com/mbrevoort/go-httpinfo


## Run Tests

    git clone https://github.com/mbrevoort/go-httpinfo && cd go-httpinfo
    go test .
