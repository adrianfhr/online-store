package repositories

import (
	"context"
	"fmt"
	"online-store/core/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	SaveTx(ctx context.Context, tx *sqlx.Tx, payment entities.Payment) (entities.Payment, error)
	GetPaymentByInvoiceID(ctx context.Context, db *sqlx.DB, invoiceID string) (entities.Payment, error)
	GetPaymentByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) ([]entities.Payment, error)
}

type PaymentRepositoryImpl struct{}

func NewPaymentRepository() PaymentRepository {
	return &PaymentRepositoryImpl{}
}

func (r *PaymentRepositoryImpl) SaveTx(ctx context.Context, tx *sqlx.Tx, payment entities.Payment) (entities.Payment, error) {

	payment.ID = uuid.New()
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	payment.Status = entities.PaymentStatusCompleted
	payment.PaymentDate = time.Now()

	query := `
        INSERT INTO payments (id, invoice_id, customer_id, amount, status, payment_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, invoice_id, customer_id, amount, status, payment_date, created_at, updated_at`

	err := tx.QueryRowContext(
		ctx, query,
		payment.ID, payment.InvoiceID, payment.CustomerID, payment.Amount, payment.Status, payment.PaymentDate, payment.CreatedAt, payment.UpdatedAt,
	).Scan(&payment.ID, &payment.InvoiceID, &payment.CustomerID, &payment.Amount, &payment.Status, &payment.PaymentDate, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		fmt.Println("Error saving payment: ", err)
		return entities.Payment{}, err
	}

	return payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByInvoiceID(ctx context.Context, db *sqlx.DB, invoiceID string) (entities.Payment, error) {
	var payment entities.Payment
	query := `SELECT * FROM payments WHERE invoice_id = $1`
	err := db.GetContext(ctx, &payment, query, invoiceID)
	if err != nil {
		return entities.Payment{}, err
	}
	return payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) ([]entities.Payment, error) {
	var payments []entities.Payment
	query := `SELECT * FROM payments WHERE customer_id = $1`
	err := db.SelectContext(ctx, &payments, query, customerID)
	if err != nil {
		return []entities.Payment{}, err
	}
	return payments, nil
}
