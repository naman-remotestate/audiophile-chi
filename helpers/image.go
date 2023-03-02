package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"gorm.io/gorm"
)

func StoreImage(filePath string, variantId uint) *gorm.DB {
	db := database.DB
	ImageInfo := models.Images{
		VariantId: variantId,
		FilePath:  filePath,
	}
	tx := db.Model(&models.Images{}).Create(&ImageInfo)
	return tx
}

func GetImagesPath(variantId uint, limit, page int64) ([]string, *gorm.DB) {
	db := database.DB
	var ImagePaths []string
	tx := db.Model(&models.Images{}).Where("variant_id = ?", variantId).Pluck("file_path", &ImagePaths).
		Limit(int(limit)).Offset(int(limit * page))
	return ImagePaths, tx
}

func UpdateImage(imageId uint, imagePath string) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Images{}).Where("id = ?", imageId).Update("file_path", imagePath)
	return tx
}
