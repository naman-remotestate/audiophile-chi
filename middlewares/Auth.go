package middlewares

import (
	"audiophile/database"
	"audiophile/models"
	"audiophile/utils"
	"context"
	"fmt"
	"net/http"
)

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if headerLength := len(r.Header["Authorization"]); headerLength == 0 {
			fmt.Fprintf(w, "auth token is not detected")
			return
		}
		authToken := r.Header["Authorization"][0]
		err := utils.ValidateToken(authToken)
		if err != nil {
			fmt.Fprintf(w, "Token is not valid,login again")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId := utils.SessionData.UserId
		sessionId := utils.SessionData.SessionId
		currentSessionInfo := models.SessionBody{}
		db := database.DB
		tx := db.Model(&models.Session{}).Where("id = ? and end_time is not null", sessionId).Find(&currentSessionInfo)
		if error := tx.Error; error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "something went wrong,please try again")
			return
		}
		if currentSessionInfo.UserId != 0 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Your session has ended, please login again")
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "UserId", userId)
		ctx = context.WithValue(ctx, "SessionId", sessionId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
