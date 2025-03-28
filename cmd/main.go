package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"prosto-delaj-api/initialize"
	"prosto-delaj-api/internal/handler"
	"prosto-delaj-api/internal/repository"
	"prosto-delaj-api/internal/service"
	"prosto-delaj-api/internal/usecase"
	"prosto-delaj-api/models"
	"prosto-delaj-api/server"
	"syscall"
)

func main() {
	appState := &models.AppState{
		Env:           &models.Environment{},
		ConfigService: &models.ConfigService{},
		ConfigVault:   &models.ConfigVault{},
	}

	logrus.Info("start init server")
	if err := initialize.RunLogger(); err != nil {
		logrus.Fatal(err.Error())
	}
	if err := initialize.LoadConfiguration(appState); err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("end init server")

	logrus.Info("start server")

	businessDatabase, queries := repository.NewBusinessDatabase(appState.ConfigService.BusinessDB)

	repo := repository.NewRepository(queries)
	svc := service.NewService(repo)
	uc := usecase.NewUsecase(svc)
	hdl := handler.NewHandler(uc)
	var serverInstance server.Server

	go runServer(&serverInstance, hdl, appState.ConfigService.Server)

	runChannelStopServer()

	serverInstance.Shutdown(context.Background(), businessDatabase)
}

func runServer(server *server.Server, hdl *handler.Handler, config *models.ServerConfig) {
	ginEngine := hdl.InitRoutes(*config)

	if err := server.Run(config.Port, ginEngine); err != nil {
		if err.Error() != "http: Server closed" {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}
}

func runChannelStopServer() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
}
