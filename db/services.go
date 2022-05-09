package db

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Delivered bool      `json:"delivered" db:"delivered"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Product struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Price  float32   `json:"price" db:"price"`
	Volume float32   `json:"volume" db:"volume"`
}

type ProductsByOrderId struct {
	ID        uuid.UUID `json:"order_items_id" db:"id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	Name      string    `json:"name" db:"name"`
	Price     float32   `json:"price" db:"price"`
	Volume    float32   `json:"volume" db:"volume"`
}

type Shipment struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Capacity  float32    `json:"capacity" db:"capacity"`
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date" db:"end_date"`
}

type ShipmentsByProductId struct {
	Capacity     float32    `json:"capacity" db:"capacity"`
	StartDate    time.Time  `json:"start_date" db:"start_date"`
	EndDate      *time.Time `json:"end_date" db:"end_date"`
	Quantity     int        `json:"quantity" db:"quantity"`
	OrderItemsId uuid.UUID  `json:"order_items_id" db:"order_items_id"`
	ShipmentId   uuid.UUID  `json:"shipment_id" db:"shipment_id"`
}

type OrderService interface {
	CreateNewOrder(order *Order) error
	AddProductToOrder(productId, orderId uuid.UUID, quantity int) error
	AddProductsToOrder(product []Product) error
}

type ProductService interface {
	GetAllProducts() ([]Product, error)
	CreateNewProduct(product *Product) error
	GetProductsByOrderId(orderId uuid.UUID) ([]ProductsByOrderId, error)
}

type ShipmentService interface {
	CreateNewShipment(shipment *Shipment) error
	AddProductToShipment(shipmentId, itemId uuid.UUID, quantity int) error
	AddProductsToShipment(product []Product) error
	GetShipmentsByItemId(itemId uuid.UUID) ([]ShipmentsByProductId, error)
	FindShipment(shipmentId uuid.UUID) (*Shipment, error)
}

type Services interface {
	OrderService
	ProductService
	ShipmentService
}
