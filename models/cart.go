package models

type Carts struct {
	Model
	UserId     uint
	User       Users `gorm:"foreignKey:UserId"`
	ItemsCount uint
	TotalCost  uint
}

type CartDetails struct {
	Model
	CartId    uint
	Cart      Carts `gorm:"foreignKey:CartId"`
	VariantId uint
	Variant   ProductVariants `gorm:"foreignKey:VariantId"`
	Quantity  uint
}
