package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ch0c0-msk/wb-tech-L0/pkg/handler"
	"github.com/ch0c0-msk/wb-tech-L0/pkg/repository"
	"github.com/ch0c0-msk/wb-tech-L0/pkg/service"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

func main() {
	postgresDB, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("ERROR: failed to initialize db: %s", err.Error())
	}
	dbRepo := repository.NewDbRepository(postgresDB)
	cacheRepo, err := repository.NewCacheRepository(postgresDB)
	if err != nil {
		log.Fatalf("ERROR: failed to create cache database - %s", err.Error())
	}
	service := service.NewService(cacheRepo, dbRepo)
	apiHandler := handler.NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("/", apiHandler.GetHomePage)
	mux.HandleFunc("/get", apiHandler.GetOrder)
	server := &http.Server{
		Addr:    ":" + viper.GetString("web.port"),
		Handler: mux,
	}

	natsHandler := handler.NewNatsHandler(service)
	conn, err := stan.Connect(viper.GetString("nats.clusterId"), viper.GetString("nats.clientId"),
		stan.NatsURL(fmt.Sprintf("nats://%s:%s", viper.GetString("nats.host"), viper.GetString("nats.port"))))
	if err != nil {
		log.Fatalf("ERROR: failed to nats connection: %s", err.Error())
	}
	defer conn.Close()

	sub, err := conn.Subscribe(viper.GetString("nats.subject"), natsHandler.AddOrder)
	if err != nil {
		log.Fatalf("ERROR: failed to nats subscription: %s", err.Error())
	}
	defer sub.Unsubscribe()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("ERROR: while starting server - %s", err.Error())
	}
}

func init() {
	initConfig()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
}

func initConfig() {
	viper.SetDefault("web.port", "8000")

	viper.SetDefault("db.username", "wb-tech-user")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5437")
	viper.SetDefault("db.dbname", "wb-tech-app")
	viper.SetDefault("db.sslmode", "disable")

	viper.SetDefault("nats.host", "localhost")
	viper.SetDefault("nats.port", "4223")
	viper.SetDefault("nats.clusterId", "test-cluster")
	viper.SetDefault("nats.cliendId", "sub")
	viper.SetDefault("nats.subject", "orderInfo")

	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("ERROR: error while reading the configuration file - %s", err.Error())
	}
}
