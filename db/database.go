package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zldobbs/ambrosia-server/graph/model"
)

// Initialize a connection pool for the Ambrosia database.
// Using pgx to connect to Postgres.
// Expects to find connection information in environment:
//   - POSTGRES_USER: Database user username
//   - POSTGRES_PASSWORD: Database user password
//   - POSTGRES_DB: Database name
//   - POSTGRES_HOST: Database hostname
//   - POSTGRES_PORT: Database port
//
// Returns:
//   - Connection pool for configured database
func InitDB() *pgxpool.Pool {
	// Load configuration from the environment
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || password == "" || dbname == "" || host == "" || port == "" {
		panic(fmt.Errorf("could not find database connection information in environment"))
	}

	// Connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	// Open a postgres connection pool
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database, error: %s", err)
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database, error: %s", err)
	}

	return pool
}

// FIXME: Remove all migration capabilities from the server code, handle manually
// FIXME: Remove all new dependencies added (golang-migrate, pq, ...)

// Run all pending SQL migrations against the database
//
// Parameters:
//   - db: raw SQL db connection
//
// Returns:
//   - Error if migration not successful
func Migrate(db *sql.DB) error {
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		return fmt.Errorf("could not find database name in environment")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		dbname,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %v", err)
	}

	// Run the migrations (up or down)
	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("migration failed: %v", err)
		} else {
			log.Println("No changes?")
		}
	}

	log.Println("Migrations applied successfully!")
	return nil
}

// Run all pending SQL migrations against the database from a pool connection.
//
// Parameters:
//   - pool: psql database pool
//
// Returns:
//   - Error if migration not successful
func MigrateFromPool(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	return Migrate(db)
}

// Construct a "WHERE" clause from a map of options.
//
// Parameters:
//   - where: map of column names and values to filter upon
//
// Returns:
//   - Tuple of query string for a WHERE statement with corresponding args in order
func BuildWhereQuery(where map[string]interface{}) (string, []interface{}) {
	if len(where) == 0 {
		return "", nil
	}

	conditions := []string{}
	args := []interface{}{}
	argPosition := 1
	for column, value := range where {
		conditions = append(conditions, column+fmt.Sprintf(" = $%d", argPosition))
		args = append(args, value)
		argPosition += 1
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

// Get a collection of ingredients from the database.
// TODO: Limit responses here, use pagination
//
// Parameters:
//   - pool: pgx databse pool connection
//   - ctx: pgx connection context
//   - where: map of filters to apply as "WHERE" clauses in the query
//
// Returns:
//   - Array of Ingredients encoded as the defined model object
func GetIngredients(pool *pgxpool.Pool, ctx context.Context, where map[string]interface{}) ([]*model.Ingredient, error) {
	query := `
		SELECT i.ingredient_id, i.name, i.description, iu.user_id, iu.name
		FROM ingredient i
		JOIN user_account iu ON i.user_id = iu.user_id
	`
	whereQuery, whereArgs := BuildWhereQuery(where)

	rows, err := pool.Query(
		ctx,
		query+whereQuery,
		whereArgs...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients froms server; error: %v", err)
	}

	var ingredients []*model.Ingredient
	for rows.Next() {
		var user model.User
		var ingredient model.Ingredient

		err := rows.Scan(
			&ingredient.IngredientID,
			&ingredient.Name,
			&ingredient.Description,
			&user.UserID,
			&user.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ingredients into struct; error: %v", err)
		}

		ingredient.User = &user
		ingredients = append(ingredients, &ingredient)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to parse through returned SQL rows; error: %v", err)
	}

	return ingredients, nil
}

// Get an ingredient from the database.
//
// Parameters:
//   - pool: pgx databse pool connection
//   - ingredient_id: ID of ingredient to retrieve
//   - ctx: pgx connection context
//
// Returns:
//   - Ingredient encoded as the defined model object
func GetIngredientById(pool *pgxpool.Pool, ingredient_id string, ctx context.Context) (*model.Ingredient, error) {
	where := map[string]interface{}{"i.ingredient_id": ingredient_id}
	ingredients, err := GetIngredients(pool, ctx, where)
	if err != nil {
		return nil, err
	}
	if len(ingredients) == 0 {
		return nil, fmt.Errorf("found no ingredients with provided id")
	}
	if len(ingredients) > 1 {
		return nil, fmt.Errorf("found multiple ingredients with provided id")
	}
	return ingredients[0], nil
}

// Get a collection of recipes from the database.
// TODO: Limit responses here, use pagination
//
// Parameters:
//   - pool: pgx databse pool connection
//   - ctx: pgx connection context
//   - where: map of filters to apply as "WHERE" clauses in the query
//
// Returns:
//   - Array of Recipes encoded as the defined model object
func GetRecipes(pool *pgxpool.Pool, ctx context.Context, where map[string]interface{}) ([]*model.Recipe, error) {
	query := `
		SELECT r.recipe_id, r.name, r.description, ru.user_id, ru.name
		FROM recipe r
		JOIN user_account ru ON r.user_id = ru.user_id
	`
	whereQuery, whereArgs := BuildWhereQuery(where)

	rows, err := pool.Query(
		ctx,
		query+whereQuery,
		whereArgs...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes froms server; error: %v", err)
	}

	var recipes []*model.Recipe
	for rows.Next() {
		var recipe model.Recipe
		var recipeUser model.User
		err := rows.Scan(
			&recipe.RecipeID,
			&recipe.Name,
			&recipe.Description,
			&recipeUser.UserID,
			&recipeUser.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not load recipe: %v", err)
		}
		recipe.User = &recipeUser

		// Get all ingredients associated with this recipe
		rows, err := pool.Query(
			ctx,
			`
			SELECT i.ingredient_id, i.name, i.description, iu.user_id, iu.name
			FROM recipe_ingredient ri
			JOIN ingredient i ON ri.ingredient_id = i.ingredient_id
			JOIN user_account iu ON i.user_id = iu.user_id
			WHERE ri.recipe_id = $1
			`,
			recipe.RecipeID,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve ingredients for recipe: %v", err)
		}

		var ingredients []*model.Ingredient
		for rows.Next() {
			var ingredient model.Ingredient
			var ingredientUser model.User
			err = rows.Scan(
				&ingredient.IngredientID,
				&ingredient.Name,
				&ingredient.Description,
				&ingredientUser.UserID,
				&ingredientUser.Name,
			)
			if err != nil {
				return nil, fmt.Errorf("could not scan out row: %v", err)
			}
			ingredient.User = &ingredientUser
			ingredients = append(ingredients, &ingredient)
		}
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("failed to parse through returned SQL rows; error: %v", err)
		}

		recipe.Ingredients = ingredients
		recipes = append(recipes, &recipe)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to parse through returned SQL rows; error: %v", err)
	}

	return recipes, nil
}

// Get a recipe from the database.
//
// Parameters:
//   - pool: pgx databse pool connection
//   - recipe_id: ID of recipe to retrieve
//   - ctx: pgx connection context
//
// Returns:
//   - Recipe encoded as the defined model object
func GetRecipeById(pool *pgxpool.Pool, recipe_id string, ctx context.Context) (*model.Recipe, error) {
	where := map[string]interface{}{"r.recipe_id": recipe_id}
	recipes, err := GetRecipes(pool, ctx, where)
	if err != nil {
		return nil, err
	}
	if len(recipes) == 0 {
		return nil, fmt.Errorf("found no recipes with provided id")
	}
	if len(recipes) > 1 {
		return nil, fmt.Errorf("found multiple recipes with provided id")
	}
	return recipes[0], nil
}
