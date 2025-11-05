CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOTNULL,
);
CREATE TABLE skus(
    id SERIAL PRIMARY KEY,
    products_id INT REFERENCES products(id) NOT NULL,
    storage INT CHECK (storage IN(8, 16, 32, 64, 128, 256, 512, 1024)),
    colour TEXT CHECK NOT NULL,
    status TEXT CHECK (status IN('in_stock', 'reserved', 'sold')),
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0
);