package handlers

import (
	"audiophile/helpers"
	"audiophile/utils"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	newLimit, newPage, err := utils.GetLimitAndPage(limit, page)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	AllProducts, tx := helpers.GetAllProducts(newLimit, newPage)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not a able to fetch products")
		return
	}
	utils.RespondJSON(w, http.StatusOK, AllProducts)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		logrus.Errorf("not able to parse the body")
		return
	}
	ProductDetail, tx := helpers.GetProductById(uint(variantId))
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to get fetch product")
		return
	}
	utils.RespondJSON(w, http.StatusOK, ProductDetail)
}
