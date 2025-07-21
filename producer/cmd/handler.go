package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/tanishashrivas/message-queue/internal"
)

func (app *application) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var orderReq internal.CreateOrderReq
	err = json.Unmarshal(body, &orderReq)
	if err != nil {
		http.Error(w, "Cannot unmarshal body", http.StatusBadRequest)
		return
	}
	orderReq.Id = uuid.New()

	order, err := json.Marshal(orderReq)
	if err != nil {
		http.Error(w, "Cannot marshal body for queueing", http.StatusInternalServerError)
		return
	}

	err = app.Redis.PushToQueue(context.Background(), "orders-queue", string(order))

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order Queued!"))
}
