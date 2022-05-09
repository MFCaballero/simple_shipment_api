package api

import (
	"github.com/MFCaballero/simple_shipment_api.git/db"
	"github.com/gorilla/mux"
)

type Handler struct {
	*mux.Router
	services db.Services
}

func NewHandler(services db.Services) *Handler {
	handler := &Handler{
		Router:   mux.NewRouter().StrictSlash(true),
		services: services,
	}

	orders := OrderHandler{store: services}
	shipments := ShipmentHandler{store: services}
	products := ProductHandler{store: services}

	handler.HandleFunc("/orders", orders.CreateNewOrderWithProducts).Methods("POST")
	handler.HandleFunc("/orders/{id}", orders.GetOrderWithProductsAndShipments).Methods("GET")
	handler.HandleFunc("/shipments", shipments.CreateNewShipmentWithItems).Methods("POST")
	handler.HandleFunc("/shipments/{id}", shipments.TrackShipment).Methods("GET")
	handler.HandleFunc("/products", products.CreateNewProduct).Methods("POST")
	handler.HandleFunc("/products", products.GetProducts).Methods("GET")

	return handler
}
