package helpers

import (
	"audiophile/database"
	"audiophile/models"
	"gorm.io/gorm"
)

type UserDetail struct {
	Email    string
	MobileNo string
}

func GetUserByEmail(body *models.Users, email string) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Users{}).Where("email = ?", email).Find(body)
	return tx
}

func GetAllUsers(limit, page int64) ([]UserDetail, *gorm.DB) {
	db := database.DB
	var AllUsers []UserDetail
	tx := db.Table("users").Select("users.email,users.mobile_no").
		Joins("inner join roles on users.id = roles.id").Where("role = ?", models.User).Find(&AllUsers).
		Limit(int(limit)).Offset(int(limit * page))
	return AllUsers, tx
}

func CreateNewUserWithRole(NewUser models.Users, role models.RolesType) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Users{}).Create(&NewUser)
	if err := tx.Error; err != nil {
		return tx
	}
	RoleInfo := models.Roles{
		UserId: NewUser.Id,
		Role:   role,
	}
	tx = db.Model(&models.Roles{}).Create(&RoleInfo)
	return tx
}

func GetUserRole(userId uint) (models.RolesType, *gorm.DB) {
	db := database.DB
	role := struct {
		Role models.RolesType
	}{}
	tx := db.Model(&models.Roles{}).Select("role").Where("user_id = ?", userId).Find(&role)
	return role.Role, tx
}

func StartSession(sessionInfo models.Session) *gorm.DB {
	db := database.DB
	tx := db.Model(&models.Session{}).Create(&sessionInfo)
	return tx
}
