package chef

type ChefMenuView struct {
	ID          int32
	Name        string
	Description string
}

type ChefItemListView struct {
	ID          int32
	Name        string
	ImageURL    string
	IsAvailable bool
}

type ChefItemDetailView struct {
	ID           int32
	Name         string
	Description  string
	Instructions string
	ImageURL     string
	Ingredients  []ChefIngredientView
	Allergens    []string
}

type ChefIngredientView struct {
	Name        string
	Measurement string
	Amount      float64
}

type ChatMessageView struct {
	ID         int32
	SenderName string
	Message    string
	SentAt     string
}
