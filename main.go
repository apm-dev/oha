package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/apm-dev/oha/src/config"
	"github.com/apm-dev/oha/src/database"
	"github.com/apm-dev/oha/src/httpserver"
	"github.com/apm-dev/oha/src/user"
	userhttp "github.com/apm-dev/oha/src/user/http"
	_userRepo "github.com/apm-dev/oha/src/user/sql"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.NewConfig()

	logLevel, err := log.ParseLevel(config.App.LogLevel)
	noError(err)
	log.SetLevel(logLevel)

	// Database
	database.ApplyDatabaseMigrations(config)
	db, err := database.NewConnection(config)
	noError(err)

	// Repository
	userRepo := _userRepo.NewUserRepo(db)

	// Service
	userSvc := user.NewService(userRepo)

	// HTTP Controller
	userCtrl := userhttp.NewUserController(userSvc)

	// HTTP Server
	server := httpserver.NewServer(userCtrl)
	server.Start(config.App.WebPort, config.App.HttpPathPrefix)

	waitForSignal()

	server.Stop()
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}

func waitForSignal() {
	var stop = make(chan struct{})
	go func() {
		var sig = make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sig)
		<-sig
		log.Info("got interrupt, shutting down...")
		stop <- struct{}{}
	}()
	<-stop
}
