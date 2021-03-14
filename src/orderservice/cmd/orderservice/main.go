package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/pkg/orderservice/transport"
	"os"
	"os/signal"
	"syscall"
)

const appID = "orderservice"

type config struct {
	ServerPort        string `envconfig:"server_port"`
	DatabaseName      string `envconfig:"database_name"`
	DatabaseAddress   string `envconfig:"database_address"`
	DatabaseUser      string `envconfig:"database_user"`
	DatabasePassword  string `envconfig:"database_password"`
	DatabaseArguments string `envconfig:"database_arguments"`
}

func main() {
	c, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	setupLogger()

	killSignalChan := getKillSignalChan()
	srv := startServer(c)

	waitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func setupLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func parseConfig() (*config, error) {
	c := config{}
	if err := envconfig.Process(appID, &c); err != nil {
		return nil, err
	}

	return &c, nil
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

func startServer(c *config) *http.Server {
	log.WithFields(log.Fields{"port": c.ServerPort}).Info("starting the server")
	db := createDbConn(c)
	router := transport.Router(db)
	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.ServerPort), Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Fatal(db.Close())
	}()

	return srv
}

func createDbConn(c *config) *sql.DB {
	arguments := c.DatabaseArguments
	if len(arguments) > 0 {
		arguments = "?" + arguments
	}

	dsn := fmt.Sprintf("%s:%s@%s/%s%s", c.DatabaseUser, c.DatabasePassword, c.DatabaseAddress, c.DatabaseName, arguments)
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
