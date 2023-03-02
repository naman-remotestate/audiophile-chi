package handlers

import (
	"audiophile/database"
	"audiophile/helpers"
	"audiophile/models"
	"audiophile/utils"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	UserInfo := models.UserBody{}
	err := utils.ParseBody(r.Body, &UserInfo)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Please check the data and try again!")
		return
	}
	Body := models.Users{}
	email := UserInfo.Email
	tx := helpers.GetUserByEmail(&Body, email)
	if error := tx.Error; error != nil {
		utils.RespondError(w, http.StatusInternalServerError, error, "something went wrong,please try again")
		return
	}

	if Body.Email != "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "User already exist,please login !")
		return
	}
	role := UserInfo.Role
	password := UserInfo.Password
	hashedPassword := utils.HashPassword(password)
	NewUser := models.Users{
		Email:    email,
		Password: hashedPassword,
		MobileNo: Body.MobileNo,
	}
	tx = helpers.CreateNewUserWithRole(NewUser, role)
	if error := tx.Error; error != nil {
		utils.RespondError(w, http.StatusInternalServerError, error, "Not able to register the user,please try again!")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusCreated, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	UserInfo := models.UserBody{}
	err := utils.ParseBody(r.Body, &UserInfo)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "please check the credentials!")
		return
	}
	email := UserInfo.Email
	User := models.Users{}
	tx := helpers.GetUserByEmail(&User, email)
	if error := tx.Error; error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Not able to login the user,please try again!")
		return
	}
	if User.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "email not registered,register first!")
		return
	}
	hashedPassword := User.Password
	password := UserInfo.Password
	userId := User.Id
	isValidUser := utils.ComparePassword(password, hashedPassword)
	if !isValidUser {
		utils.RespondError(w, http.StatusUnauthorized, nil, "wrong credentials,please check !")
		return
	}

	sessionInfo := models.Session{
		UserId:    userId,
		StartTime: time.Now(),
	}
	tx = helpers.StartSession(sessionInfo)
	if error := tx.Error; error != nil {
		utils.RespondError(w, http.StatusInternalServerError, error, "not able to start the session,please try again!")
	}
	sessionId := sessionInfo.Id
	authToken := utils.GenerateAuthToken(userId, sessionId)
	utils.RespondJSON(w, http.StatusOK, authToken)

}

func GetCartDetails(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserId").(uint)
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	newLimit, newPage, err := utils.GetLimitAndPage(limit, page)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	CartInfo, tx := helpers.GetCartByUserId(userId)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, tx.Error, "please check the entered data!")
		return
	}
	Cart, tx := helpers.GetCartDetails(CartInfo.Id, newLimit, newPage)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to fetch the cart details!")
		return
	}
	utils.RespondJSON(w, http.StatusOK, Cart)
}

func AddProductToCart(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	userId := r.Context().Value("UserId").(uint)
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	cartId, tx := helpers.CreateCartAndReturnCartId(userId)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to create the cart")
		Tx.Rollback()
		return
	}
	tx = helpers.AddProductToCart(cartId, uint(variantId))
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to add the product to cart!")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusCreated, "product added to the cart!")
}

func RemoveProductFromCart(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	userId := r.Context().Value("UserId").(uint)
	variantId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	cartId, tx := helpers.CreateCartAndReturnCartId(userId)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "something went wrong")
		Tx.Rollback()
		return
	}
	tx = helpers.RemoveFromCart(cartId, uint(variantId))
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to remove the product from cart")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusOK, "product remove from cart successfully!")
}

func GetAllAddress(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserId").(uint)
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

func AddAddress(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	userId := r.Context().Value("UserId").(uint)
	AddressInfo := models.AddressBody{}
	err := utils.ParseBody(r.Body, &AddressInfo)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, " please check the data entered!")
		return
	}
	AddressInfo.UserId = userId
	tx := helpers.StoreAddress(AddressInfo)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to add the address!")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusCreated, "successfully added the address!")
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	Tx := database.DB.Begin()
	userId := r.Context().Value("UserId").(uint)
	addressId, err := strconv.ParseInt(chi.URLParam(r, "addressId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	updatedAddressInfo := models.AddressBody{}
	err = utils.ParseBody(r.Body, &updatedAddressInfo)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, " please check the data entered!")
		return
	}
	UpdatedFields := map[string]interface{}{
		"HouseNo":     updatedAddressInfo.HouseNo,
		"StreetName":  updatedAddressInfo.StreetName,
		"CityName":    updatedAddressInfo.CityName,
		"StateName":   updatedAddressInfo.StateName,
		"CountryName": updatedAddressInfo.CountryName,
		"ZipCode":     updatedAddressInfo.ZipCode,
	}
	tx := helpers.ChangeAddressInfo(userId, uint(addressId), UpdatedFields)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to add the address!")
		Tx.Rollback()
		return
	}
	Tx.Commit()
	utils.RespondJSON(w, http.StatusOK, "address updated successfully")
}

func RemoveAddress(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserId").(uint)
	addressId, err := strconv.ParseInt(chi.URLParam(r, "addressId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	deleteInfo := map[string]interface{}{
		"DeletedAt": time.Now(),
	}
	tx := helpers.ChangeAddressInfo(userId, uint(addressId), deleteInfo)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to remove the address")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "address removed successfully")
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserId").(uint)
	orderInfo := struct {
		AddressId uint `json:"addressId"`
	}{}
	err := utils.ParseBody(r.Body, &orderInfo)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, " please check the data entered!")
		return
	}
	tx := helpers.CreateOrder(userId, orderInfo.AddressId)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to place the order")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "order placed successfully")
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserId").(uint)
	orderId, err := strconv.ParseInt(chi.URLParam(r, "orderId"), 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	tx := helpers.CancelOrder(userId, uint(orderId))
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to cancel the order")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "order cancelled successfully")
}

func AllOrders(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserId").(uint)
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	newLimit, newPage, err := utils.GetLimitAndPage(limit, page)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "not a valid request")
		return
	}
	AllOrders, tx := helpers.GetAllOrders(userId, newLimit, newPage)
	if tx.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to fetch the orders")
		return
	}
	utils.RespondJSON(w, http.StatusOK, AllOrders)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	db := database.DB
	sessionId := r.Context().Value("SessionId").(uint)
	sessionInfo := models.SessionBody{
		EndTime: time.Now(),
	}
	tx := db.Model(&models.Session{}).Where("id = ?", sessionId).Updates(sessionInfo)
	if err := tx.Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx.Error, "not able to logout")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "Successfully logged out !")
}
