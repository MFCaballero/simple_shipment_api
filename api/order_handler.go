package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MFCaballero/simple_shipment_api.git/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderHandler struct {
	store db.Services
}

func (oh *OrderHandler) CreateNewOrderWithProducts(w http.ResponseWriter, r *http.Request) {
	var products []db.OrderItems
	err := json.NewDecoder(r.Body).Decode(&products)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order := &db.Order{
		ID:        uuid.New(),
		Delivered: false,
		CreatedAt: time.Now(),
	}
	err = oh.store.CreateNewOrder(order)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{}
	response["order_status"] = fmt.Sprintf("order successfully created with id: %v", order.ID)
	for i := 0; i < len(products); i++ {
		products[i].ID = uuid.New()
		products[i].OrderId = order.ID
	}
	if err = oh.store.AddProductsToOrder(products); err != nil {
		log.Println(err)
		response["products_status"] = fmt.Sprintf("error trying to add products to order %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response["products_status"] = "success"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (oh *OrderHandler) GetOrderWithProductsAndShipments(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	orderUuid, err := uuid.Parse(orderId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	products, err := oh.store.GetProductsByOrderId(orderUuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{}
	response["products"] = products
	for _, product := range products {
		shipments, err := oh.store.GetShipmentsByItemId(product.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response[fmt.Sprintf("shipment_itemId_%s", product.ID.String())] = shipments
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
