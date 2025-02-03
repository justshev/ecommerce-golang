package storer

import (
	"context"
	"fmt"
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
tcs := []struct {
	name string
	test func(*testing.T, *MySqlStorer, sqlmock.Sqlmock)
} {
	{
		name : "success",
		test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){

			mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
			cp , err := st.CreateProduct(context.Background(), p)
			require.NoError(t, err)
			require.Equal(t, int64(1), cp.ID)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		},
	} , 
	{
		name : "failure insert product",
		test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){

			mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnError(fmt.Errorf("error"))
			_, err := st.CreateProduct(context.Background(), p)
			require.Error(t, err)

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		},
	},
	{
		name : "failure last insert id",
		test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){

			mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error")))
			_,  err := st.CreateProduct(context.Background(), p)
			require.Error(t, err)

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		},
	},
}
for _, tc := range tcs {
	t.Run(tc.name, func(t *testing.T){
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
	st := NewMySqlStorer(db)
	tc.test(t, st, mock)



})
	})
	
}



}

func TestGetProduct(t *testing.T){
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
	tcs := []struct {
		name string 
		test func(*testing.T, *MySqlStorer, sqlmock.Sqlmock)
	} {
		{
			name: "success",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
					rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
		AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
	mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)
	gp, err := st.GetProduct(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), gp.ID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
			},
		},
		{
			name: "failure get product",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error"))
				_, err := st.GetProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T){
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
	st := NewMySqlStorer(db)
	tc.test(t, st, mock)
	}) }) }
	
}

func TestListProducts(t *testing.T){
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
	tcs := []struct {
		name string 
		test func(*testing.T, *MySqlStorer, sqlmock.Sqlmock)
	} {
		{
			name: "success",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
		AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
	mock.ExpectQuery("SELECT * FROM products").WillReturnRows(rows)
	lp, err := st.ListProducts(context.Background())
	require.NoError(t, err)
	require.Len(t, lp, 1)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
			},
		},
		{
			name: "failure list products",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
				mock.ExpectQuery("SELECT * FROM products").WillReturnError(fmt.Errorf("error"))
				_, err := st.ListProducts(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T){
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
	st := NewMySqlStorer(db)
	tc.test(t, st, mock)
	}) }) }
	
}

func TestUpdateProduct(t *testing.T){
	p := &Product{
		ID: 1,
		Name: "test",
		Image: "test.jpg",
		Category: "test",
		Description: "test",
		Rating: 5,
		NumReviews: 10,
		Price: 100,
		CountInStock: 10,
	}
	np := &Product{
		ID: 1,
		Name: "test1",
		Image: "test1.jpg",
		Category: "test1",
		Description: "test1",
		Rating: 4,
		NumReviews: 11,
		Price: 101,
		CountInStock: 11,} 

	tcs := []struct {
		name string
		test func(*testing.T, *MySqlStorer, sqlmock.Sqlmock) }{
			{
				name: "success",
				test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
					mock.ExpectExec(("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")).WillReturnResult(sqlmock.NewResult(1, 1))
					cp , err := st.CreateProduct(context.Background(), p)
					require.NoError(t, err)
					require.Equal(t, int64(1), cp.ID)

					mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=? WHERE id=?").WillReturnResult(sqlmock.NewResult(1, 1))
					up, err := st.UpdateProduct(context.Background(), np)
					require.NoError(t, err)
					require.Equal(t, int64(1), up.ID)
					require.Equal(t, np.Name, up.Name)

					err = mock.ExpectationsWereMet()
					require.NoError(t, err)

			},
		},
		{
			name: "failure update product",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=? WHERE id=?").WillReturnError(fmt.Errorf("error"))
				_, err := st.UpdateProduct(context.Background(), p)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},

	
	}
	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
	st := NewMySqlStorer(db)
	tc.test(t, st, mock)
		})
	}

}

func TestDeleteProduct(t *testing.T){
	tcs := []struct {
		name string 
		test func(*testing.T, *MySqlStorer, sqlmock.Sqlmock)
	} {
		{
			name: "success",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
				mock.ExpectExec("DELETE FROM products WHERE id=?").WillReturnResult(sqlmock.NewResult(1, 1))
				err := st.DeleteProduct(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failure delete product",
			test: func(t *testing.T, st *MySqlStorer, mock sqlmock.Sqlmock){
				mock.ExpectExec("DELETE FROM products WHERE id=?").WillReturnError(fmt.Errorf("error"))
				err := st.DeleteProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock){
	st := NewMySqlStorer(db)
	tc.test(t, st, mock)
	}) }
	
}