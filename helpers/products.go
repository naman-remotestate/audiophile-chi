package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"gorm.io/gorm"
)

type result struct {
}

func GetProductVariant(conditions models.ProductVariants) (models.ProductVariants, *gorm.DB) {
	db := database.DB
	Product := models.ProductVariants{}
	tx := db.Model(&models.ProductVariants{}).Where(conditions).Find(&Product)
	return Product, tx
}

func GetProductById(variantId uint) (models.ProductVariants, *gorm.DB) {
	db := database.DB
	ProductInfo := models.ProductVariants{}
	tx := db.Model(&models.ProductVariants{}).Where("id = ?", variantId).Find(&ProductInfo)
	return ProductInfo, tx
}

func AddProductVariant(Product models.ProductVariants) (uint, *gorm.DB) {
	db := database.DB
	tx := db.Model(&models.ProductVariants{}).Create(&Product)
	return Product.Id, tx
}

func UpdateProductVariant(VariantId uint, UpdatedProduct map[string]interface{}) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.ProductVariants{}).Where("id = ?", VariantId).Updates(UpdatedProduct)
	return tx
}
func GetProductVariantsByIds(productVariantsIds []uint) ([]models.ProductVariants, *gorm.DB) {
	db := database.DB
	var ProductVariantDetail []models.ProductVariants
	tx := db.Model(&models.ProductVariants{}).Where(productVariantsIds).Find(&ProductVariantDetail)
	return ProductVariantDetail, tx
}

func GetAllProducts(limit, page int64) ([]models.ProductVariants, *gorm.DB) {
	db := database.DB
	var AllProducts []models.ProductVariants
	tx := db.Model(&models.ProductVariants{}).Find(&AllProducts).Limit(int(limit)).Offset(int(limit * page))
	return AllProducts, tx
}
func SearchProductByFilters(CompanyName, ProductName string, ProductCategory models.ProductCategories, ProductType models.ProductTypes) (models.ProductVariants, *gorm.DB) {
	db := database.DB
	conditions := models.ProductVariants{}
	if CompanyName != "" {
		conditions.CompanyName = CompanyName
	}
	if ProductName != "" {
		conditions.ProductName = ProductName
	}
	if ProductCategory != "" {
		conditions.ProductCategory = ProductCategory
	}
	if ProductName != "" {
		conditions.ProductType = ProductType
	}
	Product := models.ProductVariants{}
	tx := db.Model(&models.ProductVariants{}).Where(&conditions).Find(&Product)
	return Product, tx
}
