package models

import (
	"ambrosia-server/backend/utils"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Recipe struct {
	ID              uuid.UUID
	Name            string
	CreatorID       string
	Ingredients     []uint64
	Description     string
	Steps           []string
	PrepTimeMinutes int
	//video id
	//companion picture id list
}

func (recipeObject *Recipe) CreateRecipeObject(name string, creatorID string, ingredients []uint64, description string, steps []string, prepTimeMinutes int) Recipe {
	id := uuid.New()

	recipe := Recipe{
		ID:              id,
		Name:            name,
		CreatorID:       creatorID,
		Ingredients:     ingredients,
		Description:     description,
		Steps:           steps,
		PrepTimeMinutes: prepTimeMinutes,
	}
	return recipe
}

func (recipeObject *Recipe) SaveRecipe(db *sql.DB, recipe Recipe) error {
	ingredients, _ := utils.ToJSON(recipe.Ingredients)
	steps, _ := utils.ToJSON(recipe.Steps)

	query := `
	INSERT INTO recipes (id, name, creator_id, ingredients, description, steps, prep_time_minutes)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Exec(query, recipe.ID, recipe.Name, recipe.CreatorID, ingredients, recipe.Description, steps, recipe.PrepTimeMinutes)
	if err != nil {
		return err
	}
	return nil
}
