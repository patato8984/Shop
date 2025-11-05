CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOTNULL,
    stock INT NOTNULL,
    price INT NOTNULL,
);