package app

import (
	"os"
	"os/signal"
	"syscall"

	// _ "github.com/lib/pq"

	"github.com/Cheasezz/goTodo/internal/repository"
	"github.com/Cheasezz/goTodo/internal/service"
	"github.com/Cheasezz/goTodo/internal/transport/http"
	"github.com/Cheasezz/goTodo/pkg/postgres"
	httpserver "github.com/Cheasezz/goTodo/pkg/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	psql, err := postgres.NewPostgressDB(os.Getenv("PG_URL"))
	if err != nil {
		logrus.Fatalf("failed initialize db: %s", err.Error())
	}
	defer psql.Close()

	repos := repositories.NewRepositories(psql)
	services := service.NewServices(repos)
	handlers := http.NewHandlers(services)

	srv := httpserver.NewServer(viper.GetString("port"), handlers.Init())
	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-quit:
		logrus.Info("app - Run - signal: " + s.String())
	case err = <-srv.Notify():
		logrus.Errorf("app - Run - httpServer.Notify: %s", err)
	}

	if err := srv.Shutdown(); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	logrus.Print("TodoApp Shutting Down")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
