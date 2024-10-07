-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL, -- Values can be 'pending', 'completed', 'failed'
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    
    -- Foreign key constraints
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- Create an index on customer_id for faster lookups
CREATE INDEX idx_payments_customer_id ON payments (customer_id);

-- Create an index on invoice_id for faster lookups
CREATE INDEX idx_payments_invoice_id ON payments (invoice_id);
