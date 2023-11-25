package database

import (
	"database/sql"

	"github.com/sandronister/clean-arch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders(id,price,tax,final_price) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) List() ([]entity.Order, error) {
	rows, err := r.Db.Query("SELECT id,price,tax,final_price")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []entity.Order{}

	for rows.Next() {
		var id string
		var price, tax, final_price float64

		if err := rows.Scan(&id, &price, &tax, &final_price); err != nil {
			return nil, err
		}

		orders = append(orders, entity.Order{
			ID:         id,
			Price:      price,
			Tax:        tax,
			FinalPrice: final_price,
		})
	}

	return orders, nil

}
