package manager

type EmployeeView struct {
	ID          int32
	Name        string
	Role        string
	Email       string
	PhoneNumber string
	IsActive    bool
}

type ChatMessageView struct {
	ID         int32
	SenderName string
	Message    string
	SentAt     string
}

type RestaurantProfileView struct {
	ID             int32
	Name           string
	Address        string
	Phone          string
	Email          string
	LogoURL        string
	Description    string
	DescriptionAlt string
}

type OrderListItemView struct {
	ID          int32
	OrderNumber int32
	TableNumber int32
	WaiterName  string
	Status      string
	CreatedAt   string
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

type MenuView struct {
	ID             int32
	Name           string
	NameAlt        string
	Description    string
	DescriptionAlt string
	ImageURL       string
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
	PrepTimeMins    int32
	IsAvailable     bool
	ImageURL        string
	IsDeleted       bool
}

type AllergenCheckboxView struct {
	ID      int32
	Name    string
	Checked bool
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
