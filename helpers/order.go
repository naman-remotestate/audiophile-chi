package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"gorm.io/gorm"
	"time"
)

func CreateOrder(userId, addressId uint) *gorm.DB {
	db := database.DB
	CartInfo, tx := GetCartByUserId(userId)
	if tx.Error != nil {
		return tx
	}
	cartCost := CartInfo.TotalCost
	OrderInfo := models.Orders{
		UserId:      userId,
		AddressId:   addressId,
		TotalCost:   cartCost + ((18 * cartCost) / 100),
		OrderStatus: models.Pending,
	}
	tx = db.Model(models.Orders{}).Create(&OrderInfo)
	if tx.Error != nil {
		return tx
	}
	var CartProductsInfo []models.CartDetails
	tx = db.Model(&models.CartDetails{}).Where("cart_id = ?", CartInfo.Id).Find(&CartProductsInfo)
	if tx.Error != nil {
		return nil
	}
	Tx := database.DB.Begin()
	for _, cartProduct := range CartProductsInfo {
		variantId := cartProduct.Id
		tx := db.Model(&models.Inventory{}).Where("variant_id = ?", variantId).
			Update("stocks_left", gorm.Expr("stocks_left - ?", cartProduct.Quantity))
		if tx.Error != nil {
			return tx
		}
	}
	Tx.Commit()
	return tx
}

func CancelOrder(userId, orderId uint) *gorm.DB {
	db := database.DB
	updatedField := map[string]interface{}{
		"DeletedAt":   time.Now(),
		"OrderStatus": models.Cancelled,
	}
	tx := db.Model(&models.Orders{}).Where("id = ? and user_id", orderId, userId).Updates(updatedField)
	Cart, tx := GetCartByUserId(userId)
	if tx.Error != nil {
		return tx
	}
	var CartProductsInfo []models.CartDetails
	tx = db.Model(&models.CartDetails{}).Where("cart_id = ?", Cart.Id).Find(&CartProductsInfo)
	if tx.Error != nil {
		return nil
	}
	Tx := database.DB.Begin()
	for _, cartProduct := range CartProductsInfo {
		variantId := cartProduct.Id
		tx := db.Model(&models.Inventory{}).Where("variant_id = ?", variantId).
			Update("stocks_left", gorm.Expr("stocks_left + ?", cartProduct.Quantity))
		if tx.Error != nil {
			return tx
		}
	}
	Tx.Commit()
	return tx
}

func GetAllOrders(userId uint, limit, page int64) ([]models.Orders, *gorm.DB) {
	db := database.DB
	var AllOrders []models.Orders
	tx := db.Model(&models.Orders{}).Where("user_id = ?", userId).Find(&AllOrders).
		Limit(int(limit)).Offset(int(limit * page))
	return AllOrders, tx
}
