package models

type Inventory struct {
	Model
	VariantId    uint
	Variant      ProductVariants `gorm:"foreignKey:VariantId"`
	StockLeft    uint
	PurchaseCost uint
}
