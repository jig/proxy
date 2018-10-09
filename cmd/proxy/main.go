package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jig/proxy"
)

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	proxy.Service(os.Getenv("PROXY_SERV"), os.Getenv("PROXY_DEST"), stop)
}
