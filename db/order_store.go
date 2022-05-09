package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderStore struct {
	*sqlx.DB
}

func (os *OrderStore) CreateNewOrder(order *Order) error {
	_, err := os.Exec("insert into orders values ($1, $2, $3)", order.ID, order.Delivered, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}
	return nil
}

func (os *OrderStore) AddProductToOrder(productId, orderId uuid.UUID, quantity int) error {
	id := uuid.New()
	_, err := os.Exec("insert into order_items values($1,$2,$3,$4)", id, orderId, productId, quantity)
	if err != nil {
		return fmt.Errorf("error adding product to order: %w", err)
	}
	return nil
}

func (os *OrderStore) AddProductsToOrder(product []Product) error {
	return nil
}
