package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"water-jug-riddle-service/config"
	"water-jug-riddle-service/controller"
	"water-jug-riddle-service/service"

	"github.com/joho/godotenv"
	rice "github.com/GeertJohan/go.rice"
)

const (
	envFile = ".env"
)

func main() {
	// Define the rice box with the frontend client static files.
	appBox, err := rice.FindBox("./client/build")
	if err != nil {
		log.Fatal(err)
	}
	
	if _, err := os.Stat(envFile); !os.IsNotExist(err) {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("error loading .env file: %v", err.Error())
		}
	}

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err.Error())
	}

	svc := service.NewService()
	handler := controller.NewHandler(appBox, svc)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handler,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("error on server shutdown: %s", err.Error())
		}
	}()

	log.Printf("HTTP listener started on :%s @ %s", cfg.HTTPPort, time.Now().Format(time.RFC3339))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start http server: %s", err.Error())
	}
}
