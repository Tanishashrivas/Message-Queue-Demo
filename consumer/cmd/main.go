package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tanishashrivas/message-queue/internal"
)

type application struct {
	MQGateway *internal.MQGateway
	Redis     *internal.RedisClient
}

func main() {
	mqGateway := internal.NewMQGateway()
	rdsClient := internal.NewRedisClient()

	app := &application{
		MQGateway: mqGateway,
		Redis:     rdsClient,
	}

	go app.StartOrderWorker()
	app.ServeHTTP()
}

func (app *application) ServeHTTP() {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", os.Getenv("HOST"), os.Getenv("PORT")),
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Errorf("Failed to start the consumer", err)
	}
}
