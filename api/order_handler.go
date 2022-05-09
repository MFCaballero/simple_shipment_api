package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MFCaballero/simple_shipment_api.git/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderHandler struct {
	store db.Services
}

func (oh *OrderHandler) CreateNewOrderWithProducts(w http.ResponseWriter, r *http.Request) {
	order := &db.Order{
		ID:        uuid.New(),
		Delivered: false,
		CreatedAt: time.Now(),
	}
	errDbCreateOrder := oh.store.CreateNewOrder(order)
	if errDbCreateOrder != nil {
		log.Println(errDbCreateOrder)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var products map[uuid.UUID]int
	errDecoding := json.NewDecoder(r.Body).Decode(&products)
	if errDecoding != nil {
		log.Println(errDecoding)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Improvement: add verification for unique keys in products map
	dbErrors := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(products))
	//TODO: this works but would be more efficient using bulk sql insert instead of single
	for product, quantity := range products {
		go func(productId uuid.UUID, quantity int) {
			errDbAddProductToOrder := oh.store.AddProductToOrder(productId, order.ID, quantity)
			if errDbAddProductToOrder != nil {
				dbErrors <- errDbAddProductToOrder
				wg.Done()
				return
			}
			log.Println("done")
			dbErrors <- nil
			wg.Done()
		}(product, quantity)
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

func (oh *OrderHandler) GetOrderWithProductsAndShipments(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	orderUuid, errOrderUuid := uuid.Parse(orderId)
	if errOrderUuid != nil {
		log.Println(errOrderUuid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TODO: join everything in one query and benchmark against concurrent different calls
	products, errDbGetProductsByOrderId := oh.store.GetProductsByOrderId(orderUuid)
	if errDbGetProductsByOrderId != nil {
		log.Println(errDbGetProductsByOrderId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{}
	response["products"] = products
	for _, product := range products {
		shipments, errDbGetShipments := oh.store.GetShipmentsByItemId(product.ID)
		if errDbGetShipments != nil {
			log.Println(errDbGetShipments)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response[fmt.Sprintf("shipment_itemId_%s", product.ID.String())] = shipments
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
