package adminshared

import "strconv"

type RestaurantView struct {
	ID             int32
	Name           string
	Address        string
	Phone          string
	Email          string
	Description    string
	DescriptionAlt string
	LogoURL        string
	IsDeleted      bool
}

type FoodItemView struct {
	ID              int32
	Name            string
	NameAlt         string
	Description     string
	DescriptionAlt  string
	Instructions    string
	InstructionsAlt string
	Price           float64
	ImageURL        string
	PrepTimeMins    int32
	IsAvailable     bool
	IsDeleted       bool
	Ingredients     []ItemIngredientView
}

type AllergenCheckboxView struct {
	ID      int32
	Name    string
	Checked bool
}

type MenuView struct {
	ID             int32
	Name           string
	NameAlt        string
	Description    string
	DescriptionAlt string
	ImageURL       string
}

type MenuCheckboxView struct {
	ID      int32
	Name    string
	Checked bool
}

type IngredientView struct {
	ID          int32
	Name        string
	Description string
}

type AllergenAdminView struct {
	ID             int32
	Name           string
	NameAlt        string
	Description    string
	DescriptionAlt string
	Code           string
}

type EmployeeView struct {
	ID          int32
	Name        string
	Role        string
	Email       string
	PhoneNumber string
	IsActive    bool
}

type TableView struct {
	ID        int32
	Number    int32
	Capacity  int32
	Section   string
	Status    string
	IsDeleted bool
}

type IngredientOptionView struct {
	ID   int32
	Name string
}

type ItemIngredientView struct {
	IngredientID int32
	Name         string
	Measurement  string
	Amount       float64
}

type DashboardStats struct {
	Restaurants int64
	Menus       int64
	FoodItems   int64
	Employees   int64
	Ingredients int64
	Allergens   int64
	Tables      int64
}

type ChatMessageView struct {
	ID         int32
	SenderName string
	Message    string
	SentAt     string
}

type UsageItemView struct {
	ID   int32
	Name string
}

func Itoa(id int32) string {
	return strconv.Itoa(int(id))
}

func FormatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}
