package models

type Images struct {
	Model
	VariantId uint            `gorm:"not null"`
	Variant   ProductVariants `gorm:"foreignKey:VariantId"`
	FilePath  string          `gorm:"not null"`
}
