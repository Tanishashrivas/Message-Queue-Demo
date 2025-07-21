package internal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type CreateOrderReq struct {
	Id       uuid.UUID
	Item     string
	Quantity int
}

type OrderRes struct {
	Id         string
	StatusCode int
}

func (g *MQGateway) CreateOrder(req *CreateOrderReq) (*OrderRes, error) {
	baseUrl := os.Getenv("ProcessorBaseUrl")
	reqUrl := fmt.Sprintf("%s/orders", baseUrl)

	bodybytes, err := json.Marshal(req)
	if err != nil {
		fmt.Errorf("Error marshalling the request body", err)
	}

	gwRes, err := g.makeHttpRequest("POST", reqUrl, bodybytes)
	if err != nil {
		fmt.Errorf("Error making the request to gateway", err)
	}

	var orderRes OrderRes
	err = json.Unmarshal([]byte(gwRes), &orderRes)
	if err != nil {
		fmt.Errorf("Error unmarshalling the resonse body", err)
	}

	return &orderRes, nil
}
