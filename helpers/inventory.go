package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"gorm.io/gorm"
)

func AddNewProductToInventory(NewProduct models.Inventory) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Inventory{}).Create(&NewProduct)
	return tx
}

func UpdateProductDetailsInInventory(VariantId uint, UpdatedProductDetails map[string]interface{}) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Inventory{}).Where("variant_id = ?", VariantId).Updates(UpdatedProductDetails)
	return tx
}
