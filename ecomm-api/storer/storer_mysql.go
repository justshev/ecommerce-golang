package storer

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type MySqlStorer struct {
	db *sqlx.DB
}

func NewMySqlStorer(db *sqlx.DB) *MySqlStorer {
	return &MySqlStorer{db: db}
}

func (ms *MySqlStorer) CreateProduct(ctx context.Context, p *Product ) (*Product,error){
	res,err := ms.db.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err}
	p.ID = id
	return p, nil

}