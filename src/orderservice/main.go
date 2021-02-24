package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/transport"
	"os"
	"os/signal"
	"syscall"
)

const ServerUrl = ":8000"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	killSignalChan := getKillSignalChan()
	srv := startServer(ServerUrl)

	waitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(ch <-chan os.Signal) {
	sig := <-ch
	switch sig {
	case os.Interrupt:
		log.Info("get SIGINT")
	case syscall.SIGTERM:
		log.Info("got SIGTERM")
	}
}

func startServer(serverUrl string) *http.Server {
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	router := transport.Router()
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}
