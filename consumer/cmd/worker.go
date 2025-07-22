package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tanishashrivas/message-queue/internal"
)

func (app *application) StartOrderWorker() {
	for {
		order, err := app.Redis.ConsumeFromQueue(context.Background(), "orders-queue")
		if err != nil {
			log.Printf("Failed to consume message from the queue: %v", err)
			continue
		}

		log.Printf("Received order: %s", order)

		var orderReq internal.CreateOrderReq
		if err := json.Unmarshal([]byte(order), &orderReq); err != nil {
			log.Printf("Failed to unmarshal order JSON: %v", err)
			continue
		}

		log.Printf("Processing order with ID: %s, Item: %s, Quantity: %d",
			orderReq.Id, orderReq.Item, orderReq.Quantity)

		log.Printf("Order with ID: %s processed successfully", orderReq.Id)
	}
}
