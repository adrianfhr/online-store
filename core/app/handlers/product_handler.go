package handlers

import (
	"context"
	"fmt"
	"net/http"
	"online-store/core/domain/entities"
	"online-store/core/domain/repositories"
	"online-store/core/dto"
	"online-store/package/response"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ProductHandler struct {
	DB   *sqlx.DB
	Repo repositories.ProductRepository
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(db *sqlx.DB) *ProductHandler {
	return &ProductHandler{
		DB:   db,
		Repo: repositories.NewProductRepository(),
	}
}

// GetProducts retrieves a list of products by category
func (h *ProductHandler) GetProducts(c *gin.Context) {
	category := c.Query("category")
	print("Category: ", category)
	var products []entities.Product
	var err error

	// Jika kategori kosong, ambil semua produk
	if category == "" {
		products, err = h.Repo.GetAll(context.Background(), h.DB)
		if err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Failed to fetch all products", nil)
			return
		}
	} else {
		// Jika kategori tidak kosong, ambil produk berdasarkan kategori
		products, err = h.Repo.GetByCategory(context.Background(), h.DB, category)
		if err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Failed to fetch products by category", nil)
			return
		}
	}

	// Jika tidak ada produk yang ditemukan
	if len(products) == 0 {
		fmt.Println("No products found")
		response.RespondError(c, http.StatusNotFound, "No products found", nil)
		return
	}

	// Berikan respon sukses jika produk ditemukan
	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), products)
}

// AddProduct handles adding a new product
func (h *ProductHandler) AddProduct(c *gin.Context) {
	var addProductDTOs []dto.AddProductDTO // Use a slice to accept multiple products

	// Bind the incoming JSON to the slice of DTOs
	if err := c.ShouldBindJSON(&addProductDTOs); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	// Start a transaction
	tx, err := h.DB.Beginx()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	// Iterate over the DTOs and map them to the actual Product entity
	for _, addProductDTO := range addProductDTOs {
		product := entities.Product{
			Name:     addProductDTO.Name,
			Category: addProductDTO.Category,
			Price:    addProductDTO.Price,
			Quantity: addProductDTO.Quantity,
		}

		fmt.Println("Product: ", product)

		// Save product using transaction
		_, err := h.Repo.SaveTx(context.Background(), tx, product)
		if err != nil {
			_ = tx.Rollback()
			response.RespondError(c, http.StatusInternalServerError, "Failed to save product", nil)
			return
		}

		// Optionally, you can collect saved products for response
		// savedProducts = append(savedProducts, savedProduct)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	// Respond with a success message
	response.RespondSuccess(c, http.StatusCreated, "Products created successfully", nil)
}


// Get All Categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.Repo.GetAllCategories(context.Background(), h.DB)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Failed to get categories", nil)
		return
	}

	response.RespondSuccess(c, http.StatusOK, response.GetSuccess(), categories)
}