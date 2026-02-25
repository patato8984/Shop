CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    nickName VARCHAR(255) UNIQUE NOT NULL,
    mail VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) CHECK(role IN('user', 'admin')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);