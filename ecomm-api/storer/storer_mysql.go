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

func (ms *MySqlStorer) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product
	err := ms.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (ms *MySqlStorer) ListProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	err := ms.db.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ms *MySqlStorer) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock WHERE id=:id", p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (ms *MySqlStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}