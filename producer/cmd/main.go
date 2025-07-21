package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/tanishashrivas/message-queue/internal"
)

type application struct {
	MQGateway *internal.MQGateway
	Redis     *internal.RedisClient
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Could not load .env file: %v", err)
	}

	log.Println("Initializing MQGateway...")
	mqGateway := internal.NewMQGateway()
	log.Println("Initializing Redis client...")
	rdsClient := internal.NewRedisClient()

	app := &application{
		MQGateway: mqGateway,
		Redis:     rdsClient,
	}

	app.ServeHTTP()
}

func (app *application) ServeHTTP() {
	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	server := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start the producer: %v", err)
	}
}
