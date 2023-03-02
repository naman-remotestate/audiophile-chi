package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CartProductInfo struct {
	ProductName     string                   `json:"productName"`
	CompanyName     string                   `json:"companyName"`
	ProductCategory models.ProductCategories `json:"productCategory"`
	ProductColor    string                   `json:"productColor"`
	ProductType     models.ProductTypes      `json:"productType"`
	Cost            uint                     `json:"cost"`
	Quantity        uint                     `json:"quantity"`
}

func CreateCartAndReturnCartId(userId uint) (uint, *gorm.DB) {
	db := database.DB
	Cart := models.Carts{
		UserId:     userId,
		TotalCost:  0,
		ItemsCount: 0,
	}
	tx := db.Model(&models.Carts{}).Where("user_id = ?", userId).FirstOrCreate(&Cart)
	return Cart.Id, tx
}

func GetCartByUserId(userId uint) (models.Carts, *gorm.DB) {
	db := database.DB
	Cart := models.Carts{}
	tx := db.Model(&models.Carts{}).Where("user_id = ?", userId).Find(&Cart)
	return Cart, tx
}

func GetCartProduct(variantId uint) (models.CartDetails, *gorm.DB) {
	db := database.DB
	CartProduct := models.CartDetails{}
	tx := db.Model(&models.CartDetails{}).Where("variant_id = ?", variantId).Find(&CartProduct)
	return CartProduct, tx
}

func GetCartDetails(cartId uint, limit, page int64) ([]CartProductInfo, *gorm.DB) {
	var CartDetails []models.CartDetails
	db := database.DB
	tx := db.Model(&models.CartDetails{}).Where("cart_id = ?", cartId).Find(&CartDetails).
		Limit(int(limit)).Offset(int(limit * page))
	if err := tx.Error; err != nil {
		logrus.Errorf("not able to fetch the cart details")
		return nil, tx
	}
	var productVariantsIds []uint
	for _, cartDetail := range CartDetails {
		productVariantsIds = append(productVariantsIds, cartDetail.VariantId)
	}
	AllProductDetails, tx := GetProductVariantsByIds(productVariantsIds)
	if tx.Error != nil {
		return nil, tx
	}
	var AllCartProducts []CartProductInfo
	for i := 0; i < len(productVariantsIds); i++ {
		productDetail := CartProductInfo{
			CompanyName:     AllProductDetails[i].CompanyName,
			ProductName:     AllProductDetails[i].ProductName,
			ProductCategory: AllProductDetails[i].ProductCategory,
			ProductType:     AllProductDetails[i].ProductType,
			ProductColor:    AllProductDetails[i].ProductColor,
			Quantity:        CartDetails[i].Quantity,
			Cost:            AllProductDetails[i].SellingCost,
		}
		AllCartProducts = append(AllCartProducts, productDetail)
	}
	return AllCartProducts, tx
}

func AddProductToCart(cartId, variantId uint) *gorm.DB {
	db := database.DB
	CartProduct := models.CartDetails{
		CartId:    cartId,
		VariantId: variantId,
		Quantity:  1,
	}
	tx := db.Model(&models.CartDetails{}).
		Where("cart_id = ? and variant_id = ?", cartId, variantId).
		FirstOrCreate(&CartProduct)
	if tx.Error != nil {
		return tx
	}
	rowsAffected := tx.RowsAffected
	ProductInfo, res := GetProductById(variantId)
	if res.Error != nil {
		return res
	}
	cost := ProductInfo.SellingCost
	tx = db.Model(&models.Carts{}).Where("id = ?", cartId).
		Updates(map[string]interface{}{
			"items_count": gorm.Expr("items_count+?", 1),
			"total_cost":  gorm.Expr("total_cost+?", cost),
		})
	if tx.Error != nil || rowsAffected == 1 {
		return tx
	}
	tx = db.Model(&models.CartDetails{}).
		Where("cart_id = ? and variant_id = ?", cartId, variantId).
		Update("quantity", gorm.Expr("quantity + ?", 1))
	return tx
}

func RemoveFromCart(cartId, variantId uint) *gorm.DB {
	db := database.DB
	ProductInfo, tx := GetProductById(variantId)
	tx = db.Model(&models.CartDetails{}).
		Where("cart_id = ? and variant_id = ?", cartId, variantId).
		Update("quantity", gorm.Expr("quantity-?", 1))
	if tx.Error != nil {
		return tx
	}
	tx = db.Model(&models.Carts{}).Where("id = ?", cartId).
		Updates(map[string]interface{}{
			"items_count": gorm.Expr("items_count - ?", 1),
			"total_cost":  gorm.Expr("total_cost - ?", ProductInfo.SellingCost),
		})
	return tx
}
