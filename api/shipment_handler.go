package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	capacityFloat, errCapacityFloat := strconv.ParseFloat(capacity, 32)
	if errCapacityFloat != nil {
		log.Println(errCapacityFloat)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shipment := &db.Shipment{
		ID:        uuid.New(),
		Capacity:  float32(capacityFloat),
		StartDate: time.Now(),
	}
	errDbCreateShipment := sh.store.CreateNewShipment(shipment)
	if errDbCreateShipment != nil {
		log.Println(errDbCreateShipment)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var items map[uuid.UUID]int
	errDecoding := json.NewDecoder(r.Body).Decode(&items)
	if errDecoding != nil {
		log.Println(errDecoding)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Improvement: add verification for unique keys in products map
	dbErrors := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(items))
	//TODO: this works but would be more efficient using bulk sql insert instead of single
	for item, quantity := range items {
		go func(itemId uuid.UUID, quantity int) {
			errDbAddProductToShipment := sh.store.AddProductToShipment(shipment.ID, itemId, quantity)
			if errDbAddProductToShipment != nil {
				dbErrors <- errDbAddProductToShipment
				wg.Done()
				return
			}
			log.Println("done")
			dbErrors <- nil
			wg.Done()
		}(item, quantity)
	}
	go func() {
		wg.Wait()
		close(dbErrors)
	}()
	for err := range dbErrors {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
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
