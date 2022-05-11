package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShipmentStore struct {
	*sqlx.DB
}

func (ss *ShipmentStore) CreateNewShipment(shipment *Shipment) error {
	_, err := ss.Exec("insert into shipments values ($1, $2, $3)", shipment.ID, shipment.Capacity, shipment.StartDate)
	if err != nil {
		return fmt.Errorf("error creating shipment: %w", err)
	}
	return nil
}

func (ss *ShipmentStore) AddProductToShipment(shipmentId, itemId uuid.UUID, quantity int) error {
	_, err := ss.Exec("insert into shipment_products values($1,$2,$3)", shipmentId, itemId, quantity)
	if err != nil {
		return fmt.Errorf("error adding product to shipment: %w", err)
	}
	return nil
}

func (ss *ShipmentStore) AddProductsToShipment(items []ShipmentProducts) error {
	_, err := ss.NamedExec("insert into shipment_products (shipment_id, order_items_id, quantity) values (:shipment_id, :order_items_id, :quantity)", items)
	if err != nil {
		return fmt.Errorf("error adding products to shipment: %w", err)
	}
	return nil
}

func (ss *ShipmentStore) GetShipmentsByItemId(itemId uuid.UUID) ([]ShipmentsByProductId, error) {
	shipments := []ShipmentsByProductId{}
	if err := ss.Select(&shipments, "select shipments.capacity,shipments.start_date,shipments.end_date,quantity,order_items_id,shipment_id from shipment_products join shipments on shipments.id = shipment_products.shipment_id and shipment_products.order_items_id = $1", itemId); err != nil {
		return shipments, fmt.Errorf("error getting products by orderId: %w", err)
	}
	return shipments, nil
}

func (ss *ShipmentStore) FindShipment(shipmentId uuid.UUID) (*Shipment, error) {
	shipment := &Shipment{}
	if err := ss.Get(shipment, `SELECT * FROM shipments WHERE id = $1`, shipmentId); err != nil {
		return shipment, fmt.Errorf("error getting shipment: %w", err)
	}
	return shipment, nil
}
