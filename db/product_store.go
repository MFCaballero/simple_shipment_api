package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductStore struct {
	*sqlx.DB
}

func (ps *ProductStore) CreateNewProduct(product *Product) error {
	_, err := ps.Exec("insert into products values ($1, $2, $3, $4)", product.ID, product.Name, product.Price, product.Volume)
	if err != nil {
		return fmt.Errorf("error creating product: %w", err)
	}
	return nil
}

func (ps *ProductStore) GetAllProducts() ([]Product, error) {
	products := []Product{}
	if err := ps.Select(&products, "SELECT * FROM products"); err != nil {
		return products, fmt.Errorf("error getting users: %w", err)
	}
	return products, nil
}

func (ps *ProductStore) GetProductsByOrderId(orderId uuid.UUID) ([]ProductsByOrderId, error) {
	products := []ProductsByOrderId{}
	if err := ps.Select(&products, "select order_items.id,quantity,product_id,products.name,products.price,products.volume from order_items join products on products.id = order_items.product_id and order_items.order_id = $1", orderId); err != nil {
		return products, fmt.Errorf("error getting products by orderId: %w", err)
	}
	return products, nil
}
