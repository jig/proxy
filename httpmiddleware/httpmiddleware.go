package httpmiddleware

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jig/teereadcloser"
)

func MaxConnections(maxClients int, next http.HandlerFunc) http.HandlerFunc {
	restrictRequests := make(chan bool, maxClients)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		restrictRequests <- true
		defer func() { <-restrictRequests }()
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %-4.4s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

type CountMap struct {
	mutex sync.Mutex
	data  map[string]int
}

func Count(count *CountMap, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count.mutex.Lock()
		defer count.mutex.Unlock()
		count.data[r.URL.String()]++
		next.ServeHTTP(w, r)
	}
}

func Debug(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %-4.4s %s", r.RemoteAddr, r.Method, r.RequestURI)
		log.Printf("%+v", r)
		r.Body = ioaux.TeeReadCloser(r.Body, os.Stderr)
		next.ServeHTTP(&debugResponseWriter{w}, r)
	}
}

func httpNoOp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}

func (drw *debugResponseWriter) Write(b []byte) (int, error) {
	log.Println("...response:\n", string(b))
	return drw.ResponseWriter.Write(b)
}

type debugResponseWriter struct {
	http.ResponseWriter
}
