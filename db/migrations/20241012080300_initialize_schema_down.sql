-- 20241012080300_initialize_schema_down
-- Teardown database

DROP TABLE recipe_ingredient;
DROP TABLE recipe;
DROP TABLE ingredient;
DROP TABLE user_account;

\c postgres

DROP DATABASE ambrosia;
