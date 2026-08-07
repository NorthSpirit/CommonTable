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
	ID                   int32
	Number               int32
	Capacity             int32
	Section              string
	Status               string
	IsDeleted            bool
	OccupiedByRestaurant bool
}

type RestaurantOptionView struct {
	ID   int32
	Name string
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

type OrderListItemView struct {
	ID          int32
	OrderNumber int32
	TableNumber int32
	WaiterName  string
	Status      string
	CreatedAt   string
}

type OrderDetailView struct {
	ID          int32
	OrderNumber int32
	TableNumber int32
	WaiterName  string
	Status      string
	CreatedAt   string
	Notes       string
	Items       []OrderDetailItemView
	Payments    []OrderPaymentView
	Total       float64
	TotalPaid   float64
	Remaining   float64
}

type OrderDetailItemView struct {
	ID       int32
	Name     string
	Quantity int32
	Price    float64
	Notes    string
}

type OrderPaymentView struct {
	ID     int32
	Amount float64
	Method string
	Status string
}

type TableOptionView struct {
	ID     int32
	Number int32
}

type EmployeeOptionView struct {
	ID   int32
	Name string
}

type MenuItemOptionView struct {
	ID    int32
	Name  string
	Price float64
}

type TableStatusView struct {
	ID       int32
	Number   int32
	Capacity int32
	Occupied bool
}

func Itoa(id int32) string {
	return strconv.Itoa(int(id))
}

func FormatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}
