package repositories

import (
	"context"
	"fmt"
	"online-store/core/domain/entities"

	"github.com/jmoiron/sqlx"
)

type CustomerRepository interface {
	SaveTx(ctx context.Context, tx *sqlx.Tx, Customer entities.Customer) (entities.Customer, error)
	GetByEmail(ctx context.Context, db *sqlx.DB, email string) (entities.Customer, error)
	GetByID(ctx context.Context, db *sqlx.DB, id string) (entities.Customer, error)
}

type CustomerRepositoryImpl struct{}

func NewCustomerRepository() CustomerRepository {
	return &CustomerRepositoryImpl{}
}

// SaveTx saves a customer to the database
func (r *CustomerRepositoryImpl) SaveTx(ctx context.Context, tx *sqlx.Tx, customer entities.Customer) (entities.Customer, error) {
	query := `INSERT INTO customers (name, email, password) VALUES ($1, $2, $3) RETURNING id`

	// Menyimpan pelanggan dan mendapatkan ID pelanggan yang baru disimpan
	err := tx.QueryRowContext(ctx, query, customer.Name, customer.Email, customer.Password).Scan(&customer.ID)
	if err != nil {
		fmt.Println("Error saving customer: ", err)
		return entities.Customer{}, err
	}

	return customer, nil
}

// GetByEmail retrieves a customer by email
func (r *CustomerRepositoryImpl) GetByEmail(ctx context.Context, db *sqlx.DB, email string) (entities.Customer, error) {
	var customer entities.Customer

	query := `SELECT * FROM customers WHERE email = $1`

	err := db.GetContext(ctx, &customer, query, email)
	if err != nil {
		fmt.Println("Error fetching customer by email: ", err)
		return entities.Customer{}, err
	}

	return customer, nil
}

func (r *CustomerRepositoryImpl) GetByID(ctx context.Context, db *sqlx.DB, id string) (entities.Customer, error) {
	var customer entities.Customer

	query := `SELECT * FROM customers WHERE id = $1`

	err := db.GetContext(ctx, &customer, query, id)
	if err != nil {
		fmt.Println("Error fetching customer by id: ", err)
		return entities.Customer{}, err
	}
	return customer, nil
}