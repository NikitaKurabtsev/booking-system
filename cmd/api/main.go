package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/NikitaKurabtsev/booking-system/internal/repositories"

	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(loggerHandler)

	if err := loadConfig(); err != nil {
		logger.Error("failed to load configs", "error", err.Error())
	}

	//if err := godotenv.Load(); err != nil {
	//	logger.Error("failed to load environment variables", "error", err.Error())
	//}

	database, err := db.NewPostgresDB(db.Config{
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

	fmt.Println(database)

	respository := repositories.ResourcesPostgres{database}

	fmt.Println(respository)

	// serivce should recive logger, cache and others...

	// repository := NewRepository(database)
	//
	// TODO: new repo with db
	// TODO: new services with repo
	// TODO: new handlers with services

	for {

	}

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method == "POST" {
	//		logger.Error("Not Allowed for this Handler",
	//			slog.String("method", r.Method),
	//			slog.Time("time", time.Now()))
	//
	//		w.WriteHeader(http.StatusMethodNotAllowed)
	//		w.Write([]byte("Not Allowed"))
	//		return
	//	}
	//	logger.Info("INFO: Request received from %s",
	//		"address OK", r.RemoteAddr)
	//
	//	w.Write([]byte("Hello World"))
	//})
	//
	//logger.Info("INFO: Server starting on :8080")
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	logger.Error("ERROR: Failed to start server: %v", err)
	//}
}

func loadConfig() error {
	viper.AddConfigPath("config/db")
	viper.SetConfigName("postgres")
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
