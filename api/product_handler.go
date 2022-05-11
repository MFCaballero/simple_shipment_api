package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MFCaballero/simple_shipment_api.git/db"
	"github.com/google/uuid"
)

type ProductHandler struct {
	store db.Services
}

func (ph *ProductHandler) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	product := db.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ID = uuid.New()
	err = ph.store.CreateNewProduct(&product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ph *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ph.store.GetAllProducts()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}
