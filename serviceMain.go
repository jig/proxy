package proxy

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jig/proxy/httpmiddleware"
)

func Service(serviceAddr, destAddr string, stop <-chan os.Signal) {
	destAddrURL, err := url.Parse(destAddr)
	if err != nil {
		log.Fatalf("%q is an invalid destination", destAddr)
	}

	router := mux.NewRouter()
	router.HandleFunc("/",
		httpmiddleware.MaxConnections(64, NewProxy(512, destAddrURL).Proxy))

	server := &http.Server{
		Addr:    serviceAddr,
		Handler: router,
	}

	log.Printf("Service listening on http://%s\n", server.Addr)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-stop
	log.Println("\nService shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	log.Println("Service gracefully stopped")
}

type service struct {
	proxy *httputil.ReverseProxy
}

func NewProxy(conns int, dest *url.URL) *service {
	s := &service{
		proxy: httputil.NewSingleHostReverseProxy(dest),
	}
	// default values except were commented
	s.proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          conns, // modified to match MaxIdleConnsPerHost
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   conns, // default is 2...
	}
	return s
}

func (s *service) Proxy(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
