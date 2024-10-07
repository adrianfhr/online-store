# Online Store

## Description

This is a simple online store application built with Go, featuring a RESTful API for product management and a PostgreSQL database. It allows users to view, add, update, and delete products. The application also includes features for customer authentication and a shopping cart.

## Features

- **Product Management**: Add, update, delete, and view products.
- **Shopping Cart**: Customers can add products to a shopping cart, view cart items, and remove items.
- **Customer Authentication**: Users can register and log in to their accounts.
- **Database**: Uses PostgreSQL for data storage.
- **Docker**: Containerized application for easy setup and deployment.

## Technologies Used

- Go (Golang)
- Gin Gonic (web framework)
- PostgreSQL (database)
- Docker (containerization)
- SQLX (for database interaction)

## Prerequisites

- Docker
- Docker Compose
- Go (for local development)

## Getting Started with Docker

### Clone the Repository

```bash
git clone https://github.com/yourusername/online-store.git
cd online-store
```

### Docker Setup

Build the Docker image and run the container:

```bash
docker-compose up --build
```

### Run Migration

Run the migration to create the database schema:

```bash
docker-compose run migrate
```

### Up the Application

Start the application:

```bash
docker-compose up
```

The application will be available at `http://localhost:8081`.

## API Endpoints

- **Add Product**: `POST /Products`
  - **Body** (array of AddProductDTO):
    ```json
    [
        {
            "name": "Product 1",
            "category": "Category A",
            "price": 100.00,
            "quantity": 10
        },
        {
            "name": "Product 2",
            "category": "Category B",
            "price": 150.00,
            "quantity": 5
        }
    ]
    ```
  
- **Get Products**: `GET /Products`
  
- **Get Categories**: `GET /Products/Categories`

- **Customer Sign Up**: `POST /Customers/SignUp`
  - **Body** (CreateCustomerDTO):
    ```json
    {
        "name": "John Doe",
        "email": "john@example.com",
        "password": "securepassword"
    }
    ```

- **Customer Sign In**: `POST /Customers/SignIn`
  - **Body** (SignInCustomerDTO):
    ```json
    {
        "email": "john@example.com",
        "password": "securepassword"
    }
    ```

- **Add to Cart**: `POST /Cart/Items`
  - **Body** (AddToCartDTO):
    ```json
    {
        "product_id": "12345",
        "quantity": 2
    }
    ```

- **View Cart**: `GET /Cart/Items`

- **Remove Item from Cart**: `DELETE /Cart/Items`
  - **Body** (RemoveProductFromCartDTO):
    ```json
    {
        "product_id": "12345"
    }
    ```

- **Checkout**: `POST /Cart/Checkout`

- **Get Invoices**: `GET /Invoices`

- **Get Payments**: `GET /Payments`

- **Create Payment**: `POST /Payments`
  - **Body** (CreatePaymentDTO):
    ```json
    {
        "invoice_id": "67890"
    }
    ```
### API Collection

I also provide an API collection in Postman format located in the docs directory. You can use this collection to easily test the provided endpoints. Simply import the collection file into Postman to start exploring the API.