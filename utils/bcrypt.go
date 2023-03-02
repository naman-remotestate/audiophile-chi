package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	userPassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(userPassword, bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func ComparePassword(password string, hashedPassword string) bool {
	userPassword := []byte(password)
	passwordStoredInDb := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(passwordStoredInDb, userPassword)
	if err != nil {
		return false
	} else {
		return true
	}
}
