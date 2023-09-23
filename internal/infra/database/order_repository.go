package database

import (
	"database/sql"
	"fmt"

	"github.com/gbgomes/GoExpert/CleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindAll(page, limit int, sort string) ([]entity.Order, error) {
	offset := 0
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	var stmt string
	if page != 0 && limit != 0 {
		offset = limit * (page - 1)
		stmt = fmt.Sprintf("Select id, price, tax, final_price from orders order by id %s limit %d offset %d", sort, limit, offset)
	} else {
		stmt = fmt.Sprintf("Select id, price, tax, final_price from orders order by id %s", sort)
	}

	rows, err := r.Db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.Price, &o.Tax, &o.FinalPrice); err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
