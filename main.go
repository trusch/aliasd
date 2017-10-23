package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/trusch/aliasd/manager"
)

var storageURI = flag.String("storage", "leveldb:///srv/alias-db", "storage uri")
var listenAddr = flag.String("listen", ":80", "listen address")

func main() {
	flag.Parse()
	manager, err := manager.New(*storageURI)
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		manager.Close()
		os.Exit(0)
	}()
	log.Infof("start listening for incoming requests on %v", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, manager))
}
