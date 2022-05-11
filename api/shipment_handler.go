package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MFCaballero/simple_shipment_api.git/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ShipmentHandler struct {
	store db.Services
}

func (sh *ShipmentHandler) CreateNewShipmentWithItems(w http.ResponseWriter, r *http.Request) {
	capacity := r.URL.Query().Get("capacity")
	log.Println("capacity ", capacity)
	var items []db.ShipmentProducts
	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	capacityFloat, err := strconv.ParseFloat(capacity, 32)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shipment := &db.Shipment{
		ID:        uuid.New(),
		Capacity:  float32(capacityFloat),
		StartDate: time.Now(),
	}
	err = sh.store.CreateNewShipment(shipment)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{}
	response["shipment_status"] = fmt.Sprintf("shipment successfully created with id: %v", shipment.ID)
	for i := 0; i < len(items); i++ {
		quantity, err := sh.store.GetQuantityByOrderItemId(items[i].OrderItemsId)
		if err != nil {
			log.Println(err)
			response["items_status"] = fmt.Sprintf("error trying to get order item quantity %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		if quantity < items[i].Quantity {
			response["items_status"] = fmt.Sprintf("error bad request, quantity %v cannot be greater than %v", items[i].Quantity, quantity)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		items[i].ShipmentId = shipment.ID
	}
	if err = sh.store.AddProductsToShipment(items); err != nil {
		log.Println(err)
		response["items_status"] = fmt.Sprintf("error trying to add products to shipment %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response["items_status"] = "success"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (sh *ShipmentHandler) TrackShipment(w http.ResponseWriter, r *http.Request) {
	shipmentId := mux.Vars(r)["id"]
	shipmentUuid, errShipmentUuid := uuid.Parse(shipmentId)
	if errShipmentUuid != nil {
		log.Println(errShipmentUuid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shipment, errDbFindShipment := sh.store.FindShipment(shipmentUuid)
	if errDbFindShipment != nil {
		log.Println(errDbFindShipment)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shipment)
}
