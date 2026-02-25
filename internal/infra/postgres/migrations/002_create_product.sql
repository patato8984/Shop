CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    deleted_at TIMESTAMPTZ
);
CREATE TABLE skus(
    id SERIAL PRIMARY KEY,
    products_id INT REFERENCES products(id) ON DELETE CASCADE NOT NULL,
    storage INT CHECK (storage IN(8, 16, 32, 64, 128, 256, 512, 1024)),
    colour TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    update_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);