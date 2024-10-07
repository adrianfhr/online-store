-- Create carts table with a one-to-one relationship to customers
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_customer
      FOREIGN KEY(customer_id) 
      REFERENCES customers(id)
      ON DELETE CASCADE
);

-- Create cart_products table for many-to-many relationship between carts and products
CREATE TABLE cart_products (
    cart_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (cart_id, product_id),
    CONSTRAINT fk_cart
      FOREIGN KEY(cart_id) 
      REFERENCES carts(id)
      ON DELETE CASCADE,
    CONSTRAINT fk_product
      FOREIGN KEY(product_id) 
      REFERENCES products(id)
      ON DELETE CASCADE
);
