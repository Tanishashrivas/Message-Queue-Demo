package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tanishashrivas/message-queue/internal"
)

func (app *application) StartOrderWorker() {
	for {
		order, err := app.Redis.ConsumeFromQueue(context.Background(), "orders-queue")
		if err != nil {
			fmt.Errorf("Failed to consume message from the queue", err)
			return
		}

		fmt.Println("Order: ", order)
		orderJSON, err := json.Marshal(order)
		if err != nil {
			fmt.Errorf("Failed to marshal order to JSON", err)
			return
		}

		var orderRes internal.OrderRes
		if err := json.Unmarshal(orderJSON, &orderRes); err != nil {
			fmt.Errorf("Failed to unrmarshal order to JSON by the consumer", err)
			return
		}

		fmt.Sprintf("Order with Id: %s processed successfully", orderRes.Id)
	}
}
