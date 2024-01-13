-- +migrate Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    encryptedpassword VARCHAR(60) NOT NULL,
    isadmin BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (email)
);
