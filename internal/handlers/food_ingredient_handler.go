package handlers

import (
	"CommonTable/internal/services/foodservices"
)

type FoodIngredientHandler struct {
	FoodIngredientService *foodservices.FoodIngredientService
}
