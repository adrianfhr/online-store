package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"online-store/core/domain/entities"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

// CartRepository interface defines the methods for working with the shopping cart
type CartRepository interface {
	CreateCart(ctx context.Context, tx *sqlx.Tx, customerID string) (entities.Cart, error)
	GetCartByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) (entities.Cart, error)
	AddProductToCart(ctx context.Context, tx *sqlx.Tx, cartProduct entities.CartProduct) error
	RemoveProductFromCart(ctx context.Context, tx *sqlx.Tx, cartID string, productID string) error
	ClearCart(ctx context.Context, tx *sqlx.Tx, cartID string) error
	GetCartProductsByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) (entities.Cart, error)
}

// CartRepositoryImpl is the implementation of the CartRepository interface
type CartRepositoryImpl struct{}

// NewCartRepository creates a new instance of CartRepository
func NewCartRepository() CartRepository {
	return &CartRepositoryImpl{}
}

// CreateCart creates a new cart for a customer
func (r *CartRepositoryImpl) CreateCart(ctx context.Context, tx *sqlx.Tx, customerID string) (entities.Cart, error) {
	query := `INSERT INTO carts (customer_id) VALUES ($1) RETURNING id, customer_id`
	var cart entities.Cart
	err := tx.QueryRowContext(ctx, query, customerID).Scan(&cart.ID, &cart.CustomerID)
	if err != nil {
		fmt.Println("Error creating cart: ", err)
		return entities.Cart{}, err
	}
	return cart, nil
}

// GetCartByCustomerID retrieves a cart by the customer ID
func (r *CartRepositoryImpl) GetCartByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) (entities.Cart, error) {
	query := `SELECT id, customer_id FROM carts WHERE customer_id = $1`
	var cart entities.Cart
	err := db.GetContext(ctx, &cart, query, customerID)
	if err != nil {
		fmt.Println("Error retrieving cart: ", err)
		return entities.Cart{}, err
	}
	return cart, nil
}

// GetCartProductsByCustomerID retrieves the cart and its products by customer ID
func (r *CartRepositoryImpl) GetCartProductsByCustomerID(ctx context.Context, db *sqlx.DB, customerID string) (entities.Cart, error) {
    // Step 1: Get the cart details
    var cart entities.Cart
    cartQuery := `
        SELECT id, customer_id, created_at, updated_at
        FROM carts
        WHERE customer_id = $1`
    
    err := db.QueryRowContext(ctx, cartQuery, customerID).Scan(&cart.ID, &cart.CustomerID, &cart.CreatedAt, &cart.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            // If no cart is found, return an empty cart with the given customerID parse to uuid
            cart.CustomerID = uuid.MustParse(customerID)
            return cart, err
        }
        fmt.Println("Error retrieving cart: ", err)
        return entities.Cart{}, err
    }

    // Step 2: Get the products associated with the cart
    productsQuery := `
        SELECT cp.product_id, cp.quantity, p.name
        FROM cart_products cp
        LEFT JOIN products p ON cp.product_id = p.id
        WHERE cp.cart_id = $1`
    
    rows, err := db.QueryContext(ctx, productsQuery, cart.ID)
    if err != nil {
        fmt.Println("Error retrieving cart products: ", err)
        return entities.Cart{}, err
    }
    defer rows.Close()

    // Declare a slice to hold the products
    var products []entities.ProductInCart

    // Loop through the rows to scan the results
    for rows.Next() {
        var product entities.ProductInCart
        var quantity sql.NullInt64 // Use sql.NullInt64 to handle NULL values

        if err := rows.Scan(&product.ProductID, &quantity, &product.ProductName); err != nil {
            fmt.Println("Error scanning product row: ", err)
            return entities.Cart{}, err
        }

        // Check if quantity is valid before setting it
        if quantity.Valid {
            product.Quantity = int(quantity.Int64) // Set quantity if valid
        } else {
            product.Quantity = 0 // Set to 0 if quantity is NULL
        }

        products = append(products, product)
    }

    // Set the products in the cart
    cart.Items = products

    // Return the cart with its products
    return cart, nil
}

// AddProductToCart adds a product to the cart
func (r *CartRepositoryImpl) AddProductToCart(ctx context.Context, tx *sqlx.Tx, cartProduct entities.CartProduct) error {
    // Check if the product already exists in the cart
    var existingQuantity int
    queryCheck := `SELECT quantity FROM cart_products WHERE cart_id = $1 AND product_id = $2`
    err := tx.GetContext(ctx, &existingQuantity, queryCheck, cartProduct.CartID, cartProduct.ProductID)
    
    if err != nil {
        if err == sql.ErrNoRows {
            // Product not found in the cart, insert it
            queryInsert := `INSERT INTO cart_products (cart_id, product_id, quantity) VALUES ($1, $2, $3)`
            _, err := tx.ExecContext(ctx, queryInsert, cartProduct.CartID, cartProduct.ProductID, cartProduct.Quantity)
            if err != nil {
                fmt.Println("Error adding product to cart: ", err)
                return err
            }
        } else {
            // Some other error occurred
            fmt.Println("Error checking product in cart: ", err)
            return err
        }
    } else {
        // Product exists in the cart, update the quantity
        newQuantity := existingQuantity + cartProduct.Quantity
        queryUpdate := `UPDATE cart_products SET quantity = $1 WHERE cart_id = $2 AND product_id = $3`
        _, err := tx.ExecContext(ctx, queryUpdate, newQuantity, cartProduct.CartID, cartProduct.ProductID)
        if err != nil {
            fmt.Println("Error updating product quantity in cart: ", err)
            return err
        }
    }

    return nil
}


// RemoveProductFromCart removes a product from the cart
func (r *CartRepositoryImpl) RemoveProductFromCart(ctx context.Context, tx *sqlx.Tx, cartID string, productID string) error {
	query := `DELETE FROM cart_products WHERE cart_id = $1 AND product_id = $2`
	_, err := tx.ExecContext(ctx, query, cartID, productID)
	if err != nil {
		fmt.Println("Error removing product from cart: ", err)
		return err
	}
	return nil
}

// ClearCart removes all products from the cart
func (r *CartRepositoryImpl) ClearCart(ctx context.Context, tx *sqlx.Tx, cartID string) error {
	query := `DELETE FROM cart_products WHERE cart_id = $1`
	_, err := tx.ExecContext(ctx, query, cartID)
	if err != nil {
		fmt.Println("Error clearing cart: ", err)
		return err
	}
	return nil
}
