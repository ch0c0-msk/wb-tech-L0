package handler

import (
	"encoding/json"
	"log"

	"github.com/ch0c0-msk/wb-tech-L0/pkg/model"
	"github.com/ch0c0-msk/wb-tech-L0/pkg/service"
	"github.com/nats-io/stan.go"
)

type NatsHandler struct {
	service *service.Service
}

func NewNatsHandler(service *service.Service) *NatsHandler {
	return &NatsHandler{service: service}
}

func (nh *NatsHandler) AddOrder(msg *stan.Msg) {
	var order model.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		log.Printf("ERROR: failed to decode message data - %s", err.Error())
	}

	if err := nh.service.AddOrder(order); err != nil {
		log.Printf("ERROR: failed to add order - %s", err.Error())
		if err := nh.service.AddErrorOrder(order); err != nil {
			log.Printf("ERROR: failed to add error order in db - %s", err.Error())
		}
	}
}
