package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/tanishashrivas/message-queue/internal"
)

func (app *application) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var orderReq internal.CreateOrderReq
	err = json.Unmarshal(body, &orderReq)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		http.Error(w, "Cannot unmarshal body", http.StatusBadRequest)
		return
	}
	orderReq.Id = uuid.New()

	order, err := json.Marshal(orderReq)
	if err != nil {
		log.Printf("Error marshalling order for queueing: %v", err)
		http.Error(w, "Cannot marshal body for queueing", http.StatusInternalServerError)
		return
	}

	log.Printf("Pushing order %v to queue", orderReq.Id)
	err = app.Redis.PushToQueue(context.Background(), "orders-queue", string(order))
	if err != nil {
		log.Printf("Error pushing order to queue: %v", err)
		http.Error(w, "Failed to queue order", http.StatusInternalServerError)
		return
	}

	log.Printf("Order %v queued successfully", orderReq.Id)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order Queued!"))
}
