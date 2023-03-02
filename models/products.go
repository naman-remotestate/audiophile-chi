package models

type ProductCategories string
type ProductTypes string

const (
	HeadPhone ProductCategories = "headphone"
	EarPhone  ProductCategories = "earphone"
	Speakers  ProductCategories = "speakers"
)

const (
	Wire     ProductTypes = "wire"
	WireLess ProductTypes = "wireless"
)

type ProductVariants struct {
	Model
	ProductName     string            `gorm:"not null"`
	CompanyName     string            `gorm:"not null"`
	ProductCategory ProductCategories `gorm:"type:product_categories"`
	ProductColor    string            `gorm:"not null"`
	ProductType     ProductTypes      `gorm:"type:product_types"`
	SellingCost     uint              `gorm:"not null"`
}

type ProductImages struct {
	Model
	VariantId uint
	Variant   ProductVariants `gorm:"foreignKey:VariantId"`
	ImageUrl  string
}

type ProductDetailsBody struct {
	ProductName     string            `json:"productName"`
	CompanyName     string            `json:"companyName"`
	ProductCategory ProductCategories `json:"productCategory"`
	ProductColor    string            `json:"productColor"`
	ProductType     ProductTypes      `json:"productType"`
	PurchaseCost    uint              `json:"purchaseCost"`
	StockLeft       uint              `json:"stockLeft"`
}
