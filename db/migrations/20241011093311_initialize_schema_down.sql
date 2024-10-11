-- 20241011093311_initialize_schema_down
-- Teardown database

\c ambrosia

DROP TABLE recipe_ingredient;
DROP TABLE recipe;
DROP TABLE ingredient;
DROP TABLE user_account;

\c postgres

DROP DATABASE ambrosia;
