package storer

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)


func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	fn(db, mock)

}

func TestCreateProduct(t *testing.T){
withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
st := NewMySqlStorer(db)
p := &Product{
	Name: "test",
	Image: "test.jpg",
	Category: "test",
	Description: "test",
	Rating: 5,
	NumReviews: 10,
	Price: 100,
	CountInStock: 10,
}
mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
cp , err := st.CreateProduct(context.Background(), p)

require.NoError(t, err)
require.Equal(t, int64(1), cp.ID)
err = mock.ExpectationsWereMet()
require.NoError(t, err)
})


}

func TestGetProduct(t *testing.T){
	withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
st := NewMySqlStorer(db)
	p := &Product{
		Name: "test",
		Image: "test.jpg",
		Category: "test",
		Description: "test",
		Rating: 5,
		NumReviews: 10,
		Price: 100,
		CountInStock: 10,
	}
	rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
		AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
	mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)
	gp, err := st.GetProduct(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), gp.ID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
	})
	
}