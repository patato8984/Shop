CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    id_user INT REFERENCES users(id),
    bank_payment_id VARCHAR(255),
    addres VARCHAR(255),
    messages VARCHAR(255),
    statuse VARCHAR(255) CHECK(statuse IN('CREATED', 'WAITING_PAYMENT', 'PAID', 'FAILED', 'CANCELLED')),
    temporary_url_bank VARCHAR(255),
    total_amount INT,
    price_all_products NUMERIC(10, 2),
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE TABLE orders_items(
    id SERIAL PRIMARY KEY,
    id_order INT REFERENCES orders(id),
    id_skus INT REFERENCES skus(id),
    quantity INT,
    price_at_purchase NUMERIC(10, 2)
);