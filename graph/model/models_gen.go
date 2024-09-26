// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Ingredient struct {
	IngredientID string `json:"ingredientId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	User         *User  `json:"user"`
}

type Mutation struct {
}

type NewIngredient struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
}

type NewRecipe struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
}

type Query struct {
}

type Recipe struct {
	RecipeID    string        `json:"recipeId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Ingredients []*Ingredient `json:"ingredients"`
	User        *User         `json:"user"`
}

type User struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}
