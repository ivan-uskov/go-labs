package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/pkg/orderservice/transport"
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
	db := createDbConn()
	router := transport.Router(db)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Fatal(db.Close())
	}()

	return srv
}

func createDbConn() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@/%s", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Debugf("Connection to %s established", dsn)

	return db
}
