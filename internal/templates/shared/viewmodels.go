package shared

import "CommonTable/internal/textmarkup"

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
	Allergens    []string
}

type HomepageView struct {
	WelcomeBlocks []textmarkup.Block
	ImageOneURL   string
	ImageTwoURL   string
}

type QuickNavItemView struct {
	ID       int32
	Name     string
	ImageURL string
}
