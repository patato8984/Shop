CREATE TABLE cart (
    id_user PRIMARY KEY REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NEW(),
    update_at TIMESTAMPTZ,
    status TEXT DEFAULT 'active'
);
CREATE TABLE cart_items(
    id SERIAL PRIMARY KEY,
    id_cart INT REFERENCES cart(id_user) NOT NULL ON DELETE CASCADE,
    id_skus INT REFERENCES skus(id),
    quantity INT CHECK (quantity > 0),
    price_snapshot NUMERIC(10, 2)
);