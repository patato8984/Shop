CREATE TABLE user(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOTNULL,
    nickName VARCHAR(255) NOTNULL,
    mail VARCHAR(255) NOTNULL,
    password VARCHAR(255) NOTNULL,
    role VARCHAR(255) CHECK(role IN("user", "admin"))
);