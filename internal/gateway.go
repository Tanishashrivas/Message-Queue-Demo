package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type MQGateway struct {
	client *http.Client
}

func NewMQGateway() *MQGateway {
	return &MQGateway{
		client: &http.Client{},
	}
}

func (g *MQGateway) makeHttpRequest(method, url string, reqBody []byte) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		fmt.Errorf("Error creating the gateway request", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := g.client.Do(req)
	if err != nil {
		fmt.Errorf("Error making the request from gateway", err)
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return string(body), nil
}
