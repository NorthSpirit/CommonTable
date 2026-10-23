package shared

type BrandingView struct {
	Name    string
	LogoURL string
}

type RestaurantView struct {
	ID          int32
	Name        string
	Address     string
	Description string
	LogoURL     string
}

type MenuView struct {
	ID          int32
	Name        string
	Description string
	ImageURL    string
}

type MenuItemView struct {
	ID           int32
	Name         string
	Description  string
	Price        float64
	ImageURL     string
	PrepTimeMins int32
	HasPrepTime  bool
}

type ItemDetailView struct {
	ID           int32
	Name         string
	Description  string
	Instructions string
	ImageURL     string
	Price        float64
	Allergens    []string
	Ingredients  []string
}

type HomepageView struct {
	WelcomeHTML string
}

type QuickNavItemView struct {
	ID       int32
	Name     string
	ImageURL string
}
