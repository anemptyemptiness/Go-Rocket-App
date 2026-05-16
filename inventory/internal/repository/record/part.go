package record

import "time"

type PartType string

type Part struct {
	UUID          string    `db:"uuid"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	Price         int64     `db:"price"`
	PartType      PartType  `db:"part_type"`
	StockQuantity int64     `db:"stock_quantity"`
	CreatedAt     time.Time `db:"created_at"`
}
