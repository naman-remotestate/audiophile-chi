package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"gorm.io/gorm"
)

func StoreAddress(addressInfo models.AddressBody) *gorm.DB {
	db := database.DB
	NewAddress := models.Address{
		UserId:      addressInfo.UserId,
		HouseNo:     addressInfo.HouseNo,
		StreetName:  addressInfo.StreetName,
		CityName:    addressInfo.CityName,
		CountryName: addressInfo.CountryName,
		StateName:   addressInfo.StateName,
		ZipCode:     addressInfo.ZipCode,
	}
	tx := db.Model(&models.Address{}).Create(&NewAddress)
	return tx
}
func GetUserAddress(userId uint, limit, page int64) ([]models.Address, *gorm.DB) {
	db := database.DB
	var AllAddress []models.Address
	tx := db.Model(&models.Address{}).Where("user_id = ?", userId).Find(&AllAddress).
		Limit(int(limit)).Offset(int(limit * page))
	return AllAddress, tx
}

func ChangeAddressInfo(userId, addressId uint, updatedAddress map[string]interface{}) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Address{}).Where("user_id = ? and id = ?", userId, addressId).Updates(updatedAddress)
	return tx
}
