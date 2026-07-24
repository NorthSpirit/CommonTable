package kitchen

type KitchenOrderItemView struct {
	Name     string
	Quantity int32
	Notes    string
}

type KitchenOrderView struct {
	ID              int32
	OrderNumber     int32
	TableNumber     int32
	Status          string
	NextStatus      string
	NextStatusLabel string
	CreatedAt       string
	Notes           string
	Items           []KitchenOrderItemView
}
