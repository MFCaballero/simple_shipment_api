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
	errDecoding := json.NewDecoder(r.Body).Decode(&product)
	if errDecoding != nil {
		log.Println(errDecoding)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ID = uuid.New()
	errDbCreateProduct := ph.store.CreateNewProduct(&product)
	if errDbCreateProduct != nil {
		log.Println(errDbCreateProduct)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ph *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, errDbGetProducts := ph.store.GetAllProducts()
	if errDbGetProducts != nil {
		log.Println(errDbGetProducts)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}
