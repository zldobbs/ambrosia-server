-- Seed database with example data

\c ambrosia

-- Create some users
INSERT INTO user_account (name, password) VALUES
    ('Jeff Lebowski', 'thedude123'),
    ('Jim', 'password');

-- Create some ingredients
INSERT INTO ingredient (name, description, user_id) VALUES
    ('salt', 'common spice; table salt', 1),
    ('black pepper', 'common spice; ground black pepper', 1),
    ('raw chicken breast', 'raw, unprepared chicken breast', 2);

-- Create a recipe or two
INSERT INTO recipe (name, description, user_id) VALUES
    ('grilled chicken breast', 'Grill up some tasty chicken!', 1),
    ('oven baked chicken breast', 'Prepare this easy chicken dish in the oven', 2);

-- Link the ingredients to recipes
INSERT INTO recipe_ingredient (recipe_id, ingredient_id) VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 1),
    (2, 2);
