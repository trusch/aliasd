package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/trusch/aliasd/manager"
)

var storageURI = flag.String("storage", "leveldb:///tmp/alias-db", "storage uri")
var listenAddr = flag.String("listen", ":8080", "listen address")

func main() {
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
	log.Fatal(http.ListenAndServe(*listenAddr, manager))
}
