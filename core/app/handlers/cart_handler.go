package handlers

import (
	"context"
	"fmt"
	"net/http"
	"online-store/core/domain/entities"
	"online-store/core/domain/repositories"
	"online-store/core/dto"
	_ "online-store/package/config"
	"online-store/package/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CartHandler struct {
	DB                *sqlx.DB
	CartRepository    repositories.CartRepository
	ProductRepository repositories.ProductRepository
	InvoiceRepository repositories.InvoiceRepository
}

// NewCartHandler creates a new CartHandler
func NewCartHandler(db *sqlx.DB) *CartHandler {
	return &CartHandler{
		DB:                db,
		CartRepository:    repositories.NewCartRepository(),
		ProductRepository: repositories.NewProductRepository(),
		InvoiceRepository: repositories.NewInvoiceRepository(),
	}
}

// GetCart retrieves a cart by customer ID
func (h *CartHandler) GetCart(c *gin.Context) {
	// Get customer ID from token
	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	cart, err := h.CartRepository.GetCartByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get cart", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), cart)
}

// getcartwithproducts retrieves a cart by customer ID, including its products
func (h *CartHandler) GetCartWithProducts(c *gin.Context) {
	// Get customer ID from token
	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	cart, err := h.CartRepository.GetCartProductsByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get cart", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), cart)
}

// AddToCart adds a product to the cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	var addToCartDTO dto.AddToCartDTO
	if err := c.ShouldBindJSON(&addToCartDTO); err != nil {
		fmt.Println("Error binding JSON: ", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	cart, err := h.CartRepository.GetCartByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get cart", nil)
		return
	}

	// Check if product exists
	product, err := h.ProductRepository.GetByID(context.Background(), h.DB, addToCartDTO.ProductID)
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Product not found", nil)
		return
	}

	// add product to cart
	cartProduct := entities.CartProduct{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  addToCartDTO.Quantity,
	}

	tx, err := h.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	err = h.CartRepository.AddProductToCart(context.Background(), tx, cartProduct)
	if err != nil {
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to add product to cart", nil)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), nil)
}

//  removes a product from the cart
func (h *CartHandler) RemoveItemFromCart(c *gin.Context) {

	var removeProductFromCartDTO dto.RemoveProductFromCartDTO
	if err := c.ShouldBindJSON(&removeProductFromCartDTO); err != nil {
		fmt.Println("Error binding JSON: ", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	cart, err := h.CartRepository.GetCartByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get cart", nil)
		return
	}

	tx, err := h.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	err = h.CartRepository.RemoveProductFromCart(context.Background(), tx, cart.ID.String(), removeProductFromCartDTO.ProductID)
	if err != nil {
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to remove product from cart", nil)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), nil)
} 

func (h *CartHandler) CreateInvoice(c *gin.Context) {	
	userClaims, exists := c.Get("user")
	if !exists || userClaims == nil {
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	customerID := userClaims.(entities.Customer).ID

	// Get cart by customer ID
	cart, err := h.CartRepository.GetCartProductsByCustomerID(context.Background(), h.DB, customerID.String())
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get cart", nil)
		return
	}

	// Check if cart is empty
	if len(cart.Items) == 0 {
		response.RespondError(c, http.StatusBadRequest, "Cart is empty", nil)
		return
	}

	// Start a transaction
	tx, err := h.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	// Create invoice
	invoice := entities.Invoice{
        ID:         uuid.New(),
        CustomerID: customerID,
        Amount:     0, // This will be calculated below
        Status:     "pending",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        Expiration: time.Now().Add(60 * time.Minute), // 60 minutes
    }

	for _, item := range cart.Items {
        product, err := h.ProductRepository.GetByID(context.Background(), h.DB, item.ProductID.String())
        if err != nil {
            tx.Rollback()
			response.RespondError(c, http.StatusInternalServerError, "Failed to get stock", nil)
            return 
        }

        if product.Quantity < item.Quantity {
            tx.Rollback()
			response.RespondError(c, http.StatusBadRequest, "Insufficient stock", nil)
            return 
        }

        // Deduct the stock
        err = h.ProductRepository.UpdateStockTx(context.Background(), tx, item.ProductID.String(), product.Quantity - item.Quantity)
        if err != nil {
            tx.Rollback()
			response.RespondError(c, http.StatusInternalServerError, "Failed to update stock", nil)
            return 
        }

        // Create invoice item
        invoiceItem := entities.InvoiceItem{
            ID:          uuid.New(),
            InvoiceID:   invoice.ID,
            ProductID:   item.ProductID,
            ProductName: item.ProductName,
            Quantity:    item.Quantity,
            Price:       product.Price,
        }
        invoice.Amount += float64(item.Quantity) * invoiceItem.Price

        // Save invoice 
		fmt.Println("invoiceItem: ", invoiceItem)
		invoice.Items = append(invoice.Items, invoiceItem)
		fmt.Println("invoice: ", invoice)
		_, err = h.InvoiceRepository.SaveTx(context.Background(), tx, invoice)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error saving invoice: ", err)
			response.RespondError(c, http.StatusInternalServerError, "Failed to save invoice", nil)
			return
		}

    }

	// Clear cart
	err = h.CartRepository.ClearCart(context.Background(), tx, cart.ID.String())
	if err != nil {
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to clear cart", nil)
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	response.RespondSuccess(c, http.StatusCreated, response.GetSuccess(), invoice)
}