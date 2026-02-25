CREATE TABLE cart (
    id SERIAL PRIMARY KEY,
    id_user INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    update_at TIMESTAMPTZ,
    status TEXT DEFAULT 'active'
);
CREATE TABLE cart_items(
    id SERIAL PRIMARY KEY,
    id_cart INT REFERENCES cart(id) ON DELETE CASCADE,
    id_skus INT REFERENCES skus(id),
    quantity INT CHECK (quantity > 0),
    price_snapshot NUMERIC(10, 2)
);