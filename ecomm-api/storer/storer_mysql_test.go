package storer

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T){
mockDB,mock,err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
if err != nil {
	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
}
defer mockDB.Close()
db := sqlx.NewDb(mockDB, "sqlmock")
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

}