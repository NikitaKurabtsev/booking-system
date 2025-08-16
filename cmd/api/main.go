package main

import (
	"context"
	"github.com/NikitaKurabtsev/booking-system"
	logger2 "github.com/NikitaKurabtsev/booking-system/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaKurabtsev/booking-system/internal/handlers"
	"github.com/NikitaKurabtsev/booking-system/internal/repositories"
	"github.com/NikitaKurabtsev/booking-system/internal/services"

	_ "github.com/NikitaKurabtsev/booking-system/cmd/api/docs"
	"github.com/NikitaKurabtsev/booking-system/pkg/cache"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// @title Booking-System API
// @version 1.0
// @description API Server for Booking-System Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// TODO:
	// ensure with create booking handler
	// configure ROUTES!
	// Notification with email (MQ)
	logger := logger2.NewSLogger()

	if err := loadConfig(); err != nil {
		logger.Error("failed to load configs", "error", err.Error())
	}

	postgresDB, err := db.NewPostgresDB(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("DB_PASSWORD"),
	})
	if err != nil {
		logger.Error("failed to connect to the database", "error", err.Error())
	}
	logger.Info("successfully connected to the database")

	redisCache, err := cache.NewCache("redis:6379")
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
		return
	}
	logger.Info("successfully connected to the redis")

	repository := repositories.NewRepository(postgresDB, redisCache)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service, logger)

	server := new(booking.Server)
	go func() {
		if err := server.Run("8080", handler.InitRoutes()); err != nil {
			log.Fatalf("error occurred while running http server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := server.ShutDown(context.Background()); err != nil {
		logger.Error("error occurred on server shutting down", "error", err)
	}

	postgresDB.Close()
}

func loadConfig() error {
	viper.AddConfigPath("config/db")
	viper.SetConfigName("postgres")
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
