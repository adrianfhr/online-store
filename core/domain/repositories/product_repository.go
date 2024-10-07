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
	GetByID(ctx context.Context, db *sqlx.DB, id string) (entities.Product, error)
	GetStockByID(ctx context.Context, db *sqlx.DB, id string) (int, error)
	UpdateStockTx(ctx context.Context, tx *sqlx.Tx, id string, quantity int) error
	GetAllCategories(ctx context.Context, db *sqlx.DB) ([]string, error)
}

type ProductRepositoryImpl struct{}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

// SaveTx saves a product to the database
func (r *ProductRepositoryImpl) SaveTx(ctx context.Context, tx *sqlx.Tx, product entities.Product) (entities.Product, error) {
	query := `INSERT INTO products (name, category, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id`

	// Menyimpan produk dan mendapatkan ID produk yang baru disimpan
	err := tx.QueryRowContext(ctx, query, product.Name, product.Category, product.Price, product.Quantity).Scan(&product.ID)
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

func (r *ProductRepositoryImpl) GetByID(ctx context.Context, db *sqlx.DB, id string) (entities.Product, error) {
	var product entities.Product

	query := `SELECT * FROM products WHERE id = $1`

	err := db.GetContext(ctx, &product, query, id)
	if err != nil {
		fmt.Println("Error fetching product by id: ", err)
		return entities.Product{}, err
	}

	return product, nil
}

func (r *ProductRepositoryImpl) GetStockByID(ctx context.Context, db *sqlx.DB, id string) (int, error) {
	var stock int

	query := `SELECT quantity FROM products WHERE id = $1`

	err := db.GetContext(ctx, &stock, query, id)
	if err != nil {
		fmt.Println("Error fetching stock by id: ", err)
		return 0, err
	}

	return stock, nil
}

func (r *ProductRepositoryImpl) UpdateStockTx(ctx context.Context, tx *sqlx.Tx, id string, quantity int) error {
	query := `UPDATE products SET quantity = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, quantity, id)
	if err != nil {
		fmt.Println("Error updating stock: ", err)
		return err
	}

	return nil
}

func (r *ProductRepositoryImpl) GetAllCategories(ctx context.Context, db *sqlx.DB) ([]string, error) {
	var categories []string

	query := `SELECT DISTINCT category FROM products`

	err := db.SelectContext(ctx, &categories, query)
	if err != nil {
		fmt.Println("Error fetching all categories: ", err)
		return nil, err
	}

	return categories, nil
}