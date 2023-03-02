package middlewares

import (
	"audiophile/helpers"
	"audiophile/models"
	"audiophile/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := utils.SessionData.UserId
		role, tx := helpers.GetUserRole(userId)
		if tx.Error != nil {
			logrus.Errorf("not able to verify the user as admin ,%v", tx.Error)
			return
		}
		if role != models.Admin {
			return
		}
		next.ServeHTTP(w, r)
	})
}
