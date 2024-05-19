-- Initialize the ambrosia database

-- Create and connect to the database

CREATE DATABASE ambrosia;
\c ambrosia

-- Create tables for use within app

-- TODO: Need better authentication for users
CREATE TABLE "user" (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    password TEXT
);

CREATE TABLE recipe (
    recipe_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE,
    name VARCHAR(255),
    description VARCHAR(255)
);

CREATE TABLE ingredient (
    ingredient_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE,
    name VARCHAR(255),
    description VARCHAR(255)
);

CREATE TABLE recipe_ingredient (
    recipe_id INT REFERENCES recipe (recipe_id) ON UPDATE CASCADE ON DELETE CASCADE,
    ingredient_id INT REFERENCES ingredient (ingredient_id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT recipe_ingredient_id PRIMARY KEY (recipe_id, ingredient_id)
);
