package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"online-store/core/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// InvoiceRepository interface
type InvoiceRepository interface {
	SaveTx(ctx context.Context, tx *sqlx.Tx, invoice entities.Invoice) (entities.Invoice, error)
	GetInvoiceByID(ctx context.Context,db *sqlx.DB, id string) (entities.Invoice, error)
	GetInvoicesByCustomerID(ctx context.Context,db *sqlx.DB, id string) ([]entities.Invoice, error)
	FinishInvoice(ctx context.Context, tx *sqlx.Tx, id string, payment_id string) error
}

// InvoiceRepositoryImpl struct
type InvoiceRepositoryImpl struct{}

// NewInvoiceRepository creates a new instance of InvoiceRepository
func NewInvoiceRepository() InvoiceRepository {
	return &InvoiceRepositoryImpl{}
}

// SaveTx saves the invoice along with its items within a transaction
func (r *InvoiceRepositoryImpl) SaveTx(ctx context.Context, tx *sqlx.Tx, invoice entities.Invoice) (entities.Invoice, error) {
	// Generate a new UUID for the invoice
	invoice.ID = uuid.New()

	// Set the created and updated timestamps
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	// SQL query to insert the invoice into the database

	fmt.Println("invoice : ", invoice)

	query := `
		INSERT INTO invoices (id, customer_id, amount, status, created_at, updated_at, expiration)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, customer_id, amount, status, created_at, updated_at, expiration`

	// Execute the query
	if err := tx.QueryRowContext(ctx, query,
		invoice.ID, invoice.CustomerID, invoice.Amount,
		invoice.Status, invoice.CreatedAt, invoice.UpdatedAt, invoice.Expiration,
	).Scan(&invoice.ID, &invoice.CustomerID, &invoice.Amount,
		&invoice.Status, &invoice.CreatedAt, &invoice.UpdatedAt, &invoice.Expiration); err != nil {
		return entities.Invoice{}, err
	}

	// Insert invoice items
	for _, item := range invoice.Items {
		item.ID = uuid.New() // Generate a new UUID for the invoice item

		// SQL query to insert the invoice item
		itemQuery := `
			INSERT INTO invoice_items (id, invoice_id, product_id, product_name, quantity, price)
			VALUES ($1, $2, $3, $4, $5, $6)`

		// Execute the item query
		if _, err := tx.ExecContext(ctx, itemQuery,
			item.ID, invoice.ID, item.ProductID, item.ProductName, item.Quantity, item.Price,
		); err != nil {
			return entities.Invoice{}, err
		}
	}

	return invoice, nil
}

// GetInvoiceByID retrieves an invoice by ID along with its items
func (r *InvoiceRepositoryImpl) GetInvoiceByID(ctx context.Context,db *sqlx.DB, id string) (entities.Invoice, error) {
    // Step 1: Get the invoice details
    var invoice entities.Invoice
    invoiceQuery := `
        SELECT id,  customer_id, amount, status, created_at, updated_at, expiration
        FROM invoices
        WHERE id = $1`
    
    err := db.QueryRowContext(ctx, invoiceQuery, id).Scan(&invoice.ID, &invoice.CustomerID, &invoice.Amount, &invoice.Status, &invoice.CreatedAt, &invoice.UpdatedAt, &invoice.Expiration)
    if err != nil {
        if err == sql.ErrNoRows {
            // If no invoice is found, return an empty invoice
            return entities.Invoice{}, err
        }
        fmt.Println("Error retrieving invoice: ", err)
        return entities.Invoice{}, err
    }

    // Step 2: Get the items associated with the invoice
    itemsQuery := `
        SELECT id, invoice_id, product_id, product_name, quantity, price
        FROM invoice_items
        WHERE invoice_id = $1`
    
    rows, err := db.QueryContext(ctx, itemsQuery, invoice.ID)
    if err != nil {
        fmt.Println("Error retrieving invoice items: ", err)
        return entities.Invoice{}, err
    }
    defer rows.Close()

    // Declare a slice to hold the invoice items
    var items []entities.InvoiceItem

    // Loop through the rows to scan the results
    for rows.Next() {
        var item entities.InvoiceItem

        if err := rows.Scan(&item.ID, &item.InvoiceID, &item.ProductID, &item.ProductName, &item.Quantity, &item.Price); err != nil {
            fmt.Println("Error scanning invoice item row: ", err)
            return entities.Invoice{}, err
        }

        items = append(items, item)
    }

    // Set the items in the invoice
    invoice.Items = items

    // Return the invoice with its items
    return invoice, nil
}

// GetInvoicesByCustomerID retrieves all invoices for a given customer ID along with their items
func (r *InvoiceRepositoryImpl) GetInvoicesByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) ([]entities.Invoice, error) {
    // Step 1: Get the invoices associated with the customer
    invoicesQuery := `
        SELECT id, customer_id, amount, status, created_at, updated_at, expiration
        FROM invoices
        WHERE customer_id = $1`
    
    rows, err := db.QueryContext(ctx, invoicesQuery, customerID)
    if err != nil {
        fmt.Println("Error retrieving invoices: ", err)
        return nil, err
    }
    defer rows.Close()

    // Declare a slice to hold the invoices
    var invoices []entities.Invoice

    // Loop through the rows to scan the results
    for rows.Next() {
        var invoice entities.Invoice

        if err := rows.Scan(&invoice.ID, &invoice.CustomerID, &invoice.Amount, &invoice.Status, &invoice.CreatedAt, &invoice.UpdatedAt, &invoice.Expiration); err != nil {
            fmt.Println("Error scanning invoice row: ", err)
            return nil, err
        }

        // Step 2: Get the items associated with the invoice
        itemsQuery := `
            SELECT id, invoice_id, product_id, product_name, quantity, price
            FROM invoice_items
            WHERE invoice_id = $1`
        
        itemRows, err := db.QueryContext(ctx, itemsQuery, invoice.ID)
        if err != nil {
            fmt.Println("Error retrieving invoice items: ", err)
            return nil, err
        }
        defer itemRows.Close()

        // Declare a slice to hold the invoice items
        var items []entities.InvoiceItem

        // Loop through the rows to scan the invoice items
        for itemRows.Next() {
            var item entities.InvoiceItem

            if err := itemRows.Scan(&item.ID, &item.InvoiceID, &item.ProductID, &item.ProductName, &item.Quantity, &item.Price); err != nil {
                fmt.Println("Error scanning invoice item row: ", err)
                return nil, err
            }

            items = append(items, item)
        }

        // Set the items in the invoice
        invoice.Items = items

        // Append the invoice to the slice
        invoices = append(invoices, invoice)
    }

    // Return the invoices with their items
    return invoices, nil
}

// FinishInvoice updates the status of an invoice to "completed"
func (r *InvoiceRepositoryImpl) FinishInvoice(ctx context.Context, tx *sqlx.Tx, id string, payment_id string) error {
	// Update the status of the invoice to "completed"
	updateQuery := `
		UPDATE invoices
		SET status = $1, payment_id = $2
		WHERE id = $3`
	
	_, err := tx.ExecContext(ctx, updateQuery, entities.InvoiceStatusPaid, payment_id, id)
	if err != nil {
		fmt.Println("Error updating invoice status: ", err)
		return err
	}

	return nil
}