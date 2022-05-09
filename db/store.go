package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	*OrderStore
	*ProductStore
	*ShipmentStore
}

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Store{
		OrderStore:    &OrderStore{DB: db},
		ProductStore:  &ProductStore{DB: db},
		ShipmentStore: &ShipmentStore{DB: db},
	}, nil
}
