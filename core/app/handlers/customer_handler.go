package handlers

import (
	"context"
	"fmt"
	"net/http"
	"online-store/core/domain/entities"
	"online-store/core/domain/repositories"
	"online-store/core/dto"
	"online-store/package/response"
	"online-store/package/config"


	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

)

type CustomerHandler struct {
	DB   *sqlx.DB
	CustomerRepository repositories.CustomerRepository
	CartRepository repositories.CartRepository
}

// NewCustomerHandler creates a new CustomerHandler
func NewCustomerHandler(db *sqlx.DB) *CustomerHandler {
	return &CustomerHandler{
		DB:   db,
		CustomerRepository: repositories.NewCustomerRepository(),
		CartRepository: repositories.NewCartRepository(),
	}
}

// CreateCustomer creates a new customer (Sign Up)
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var createCustomerDTO dto.CreateCustomerDTO
	if err := c.ShouldBindJSON(&createCustomerDTO); err != nil {
		fmt.Println("Error binding JSON: ", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request",  nil)
		return 
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createCustomerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password: ", err)
		response.RespondError(c, http.StatusInternalServerError, "Failed to hash password", nil)
		return 
	}

	customer := entities.Customer{
		Name:  createCustomerDTO.Name,
		Email: createCustomerDTO.Email,
		Password: string(hashedPassword),
	}
	
	// Start a transaction
	tx, err := h.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		fmt.Println("Error starting transaction: ", err)
		response.RespondError(c, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	// Check if the customer already exists
	_, err = h.CustomerRepository.GetByEmail(context.Background(), h.DB, customer.Email)
	if err == nil {
		fmt.Println("Customer already exists")
		tx.Rollback()
		response.RespondError(c, http.StatusConflict, "Customer already exists", nil)
		return
	}

	fmt.Println("customer: ", customer)
	
	// Save the customer
	customer, err = h.CustomerRepository.SaveTx(context.Background(), tx, customer)
	if err != nil {
		fmt.Println("Error saving customer: ", err)
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to save customer", nil)
		return
	}

	fmt.Println("Customer ID: ", customer.ID)
	_, err = h.CartRepository.CreateCart(context.Background(), tx, customer.ID.String())
	if err != nil {
		fmt.Println("Error creating cart: ", err)
		tx.Rollback()
		response.RespondError(c, http.StatusInternalServerError, "Failed to create cart", nil)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction: ", err)
		response.RespondError(c, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	// Remove the password from the response
	customer.Password = ""
	response.RespondSuccess(c, http.StatusCreated, "Customer created successfully", customer)

}

// SignIn handles customer sign in
func (h *CustomerHandler) SignInCustomer(c *gin.Context) {
	var signInDTO dto.SignInCustomerDTO
	if err := c.ShouldBindJSON(&signInDTO); err != nil {
		fmt.Println("Error binding JSON: ", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	// Get the customer by email
	customer, err := h.CustomerRepository.GetByEmail(context.Background(), h.DB, signInDTO.Email)
	if err != nil {
		fmt.Println("Error getting customer by email: ", err)
		response.RespondError(c, http.StatusInternalServerError, "Failed to get customer", nil)
		return
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(signInDTO.Password))
	if err != nil {
		fmt.Println("Error comparing password: ", err)
		response.RespondError(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Genereate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": customer.ID,
		"exp": jwt.TimeFunc().AddDate(0, 0, 1).Unix(),
	})

	cfg := config.GetConfig()


	tokenString , err := token.SignedString([]byte(cfg.JWTSecret))

	if err != nil {
		fmt.Println("Error signing token: ", err)
		response.RespondError(c, http.StatusInternalServerError, "Failed to sign token", nil)
		return
	}

	// give the token in the response
	response.RespondSuccess(c, http.StatusOK, "Sign in successful", gin.H{ "token": tokenString})
}