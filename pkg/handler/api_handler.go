package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ch0c0-msk/wb-tech-L0/pkg/service"
)

type ApiHandler struct {
	service *service.Service
}

func (h *ApiHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/template/index.html")
}

func (h *ApiHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Query().Get("orderID")
	order, err := h.service.GetOrder(orderId)
	if err != nil {
		log.Printf("ERROR: failed to get order - %s", err.Error())
		w.Write([]byte("error: failed to get order by this id - " + orderId))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		log.Printf("ERROR: failed marshal order info - %s", err.Error())
		w.Write([]byte("error: failed to encode order info"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

func NewHandler(service *service.Service) *ApiHandler {
	return &ApiHandler{service: service}
}
