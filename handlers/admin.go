package handlers

import (
	"audiophile/database"
	"audiophile/helpers"
	"audiophile/models"
	"audiophile/utils"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func AddToInventory(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	ProductDetails := models.ProductDetailsBody{}
	err := utils.ParseBody(r.Body, &ProductDetails)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "please check the data and try again")
		return
	}
	searchConditions := models.ProductVariants{
		ProductName:     ProductDetails.ProductName,
		CompanyName:     ProductDetails.CompanyName,
		ProductCategory: ProductDetails.ProductCategory,
		ProductType:     ProductDetails.ProductType,
		ProductColor:    ProductDetails.ProductColor,
	}
	Product, tx := helpers.GetProductVariant(searchConditions)
	fmt.Println(tx.Error)
	if err := tx.Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "something went wrong,try again")
		return
	}
	fmt.Println(Product)
	if Product.Id != 0 {
		utils.RespondError(w, http.StatusBadRequest, nil, "product already exists")
		return
	}
	NewProduct := models.ProductVariants{
		ProductName:     ProductDetails.ProductName,
		CompanyName:     ProductDetails.CompanyName,
		ProductCategory: ProductDetails.ProductCategory,
		ProductType:     ProductDetails.ProductType,
		ProductColor:    ProductDetails.ProductColor,
		SellingCost:     ProductDetails.PurchaseCost + ((12 * ProductDetails.PurchaseCost) / 100),
	}
	id, tx := helpers.AddProductVariant(NewProduct)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to add the product to inventory")
		return
	}
	ProductInventoryDetails := models.Inventory{
		VariantId:    id,
		StockLeft:    ProductDetails.StockLeft,
		PurchaseCost: ProductDetails.PurchaseCost,
	}
	tx = helpers.AddNewProductToInventory(ProductInventoryDetails)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to add the product to inventory")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusCreated, "product added successfully")
}

func UpdateProductInfo(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "something went wrong,try again")
		return
	}
	UpdatedProductDetails := struct {
		StockLeft    uint `json:"stockLeft"`
		PurchaseCost uint `json:"purchaseCost"`
	}{}
	err = utils.ParseBody(r.Body, &UpdatedProductDetails)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "please check the data and try again")
		return
	}
	ProductDetails := map[string]interface{}{
		"StockLeft":    UpdatedProductDetails.StockLeft,
		"PurchaseCost": UpdatedProductDetails.PurchaseCost,
	}
	tx := helpers.UpdateProductDetailsInInventory(uint(variantId), ProductDetails)
	if err = tx.Error; err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not able to update the product")
		return
	}
	ProductVariantInfo := map[string]interface{}{
		"SellingCost": UpdatedProductDetails.PurchaseCost + ((UpdatedProductDetails.PurchaseCost * 12) / 100),
	}
	tx = helpers.UpdateProductVariant(uint(variantId), ProductVariantInfo)
	if err = tx.Error; err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not able to update the product")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusOK, "successfully updated the product variant details")
}

func DeleteProductFromInventory(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		logrus.Errorf("not able to parse the body")
		return
	}
	deleteInfo := map[string]interface{}{
		"DeletedAt": time.Now(),
	}
	tx := helpers.UpdateProductDetailsInInventory(uint(variantId), deleteInfo)
	if err = tx.Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to delete the product")
		return
	}
	tx = helpers.UpdateProductVariant(uint(variantId), deleteInfo)
	if err = tx.Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to delete the product")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusOK, "product deleted successfully")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	newLimit, newPage, err := utils.GetLimitAndPage(limit, page)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	AllUsers, tx := helpers.GetAllUsers(newLimit, newPage)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to fetch the users")
		return
	}
	utils.RespondJSON(w, http.StatusOK, AllUsers)
}

func GetAllAddressOfUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "userID is not valid")
		return
	}
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	newLimit, newPage, err := utils.GetLimitAndPage(limit, page)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	AllAddress, tx := helpers.GetUserAddress(uint(userId), newLimit, newPage)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to fetch the address")
		return
	}
	utils.RespondJSON(w, http.StatusOK, AllAddress)
}

func UploadProductImage(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	file, fileHeader, err := r.FormFile("image")
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "please upload the valid image")
		return
	}
	imagePath, err := utils.UploadImage(file, fileHeader)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "not able to upload the image")
		return
	}
	tx := helpers.StoreImage(imagePath, uint(variantId))
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to upload the image")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusCreated, "image uploaded successfully!")
}

func GetProductImages(w http.ResponseWriter, r *http.Request) {
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	newLimit, newPage, err := utils.GetLimitAndPage(limit, page)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}

	AllImages, err := utils.GetImageUrl(uint(variantId), newLimit, newPage)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "not able to fetch the images")
		return
	}
	utils.RespondJSON(w, http.StatusOK, AllImages)

}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	imageId, err := strconv.ParseInt(chi.URLParam(r, "imageId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	file, fileHeader, err := r.FormFile("image")
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "please upload the valid image")
		return
	}
	imagePath, err := utils.UploadImage(file, fileHeader)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "not able to upload the image")
		return
	}
	tx := helpers.UpdateImage(uint(imageId), imagePath)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to update the image")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusOK, "image updated successfully")
}
