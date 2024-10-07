package repositories

import (
	"context"
	"fmt"
	"online-store/core/domain/entities"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	SaveTx(ctx context.Context, tx *sqlx.Tx, product entities.Product) (entities.Product, error)
	GetAll(ctx context.Context, db *sqlx.DB) ([]entities.Product, error)
	GetByCategory(ctx context.Context, db *sqlx.DB, category string) ([]entities.Product, error)
}

type ProductRepositoryImpl struct{}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

// SaveTx saves a product to the database
func (r *ProductRepositoryImpl) SaveTx(ctx context.Context, tx *sqlx.Tx, product entities.Product) (entities.Product, error) {
	query := `INSERT INTO products (name, category, price) VALUES ($1, $2, $3) RETURNING id`

	// Menyimpan produk dan mendapatkan ID produk yang baru disimpan
	err := tx.QueryRowContext(ctx, query, product.Name, product.Category, product.Price).Scan(&product.ID)
	if err != nil {
		fmt.Println("Error saving product: ", err)
		return entities.Product{}, err
	}

	return product, nil
}

// GetAll retrieves all products
func (r *ProductRepositoryImpl) GetAll(ctx context.Context, db *sqlx.DB) ([]entities.Product, error) {
	var products []entities.Product

	query := `SELECT * FROM products`

	err := db.SelectContext(ctx, &products, query)
	if err != nil {
		fmt.Println("Error fetching all products: ", err)
		return nil, err
	}

	return products, nil
}

// GetByCategory retrieves products by category
func (r *ProductRepositoryImpl) GetByCategory(ctx context.Context, db *sqlx.DB, category string) ([]entities.Product, error) {
	var products []entities.Product

	query := `SELECT * FROM products WHERE category = $1`

	err := db.SelectContext(ctx, &products, query, category)
	if err != nil {
		fmt.Println("Error fetching products by category: ", err)
		return nil, err
	}

	return products, nil
}
