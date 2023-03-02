package models

import (
	"gorm.io/gorm"
	"time"
)

type RolesType string

const (
	Admin    RolesType = "admin"
	SubAdmin RolesType = "subadmin"
	User     RolesType = "user"
)

type Model struct {
	Id        uint `gorm:"primaryKey;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Users struct {
	Model
	Email      string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	MobileNo   string    `gorm:"not null"`
	IsArchived bool      `gorm:"default:false"`
	Address    []Address `gorm:"foreignKey:UserId"`
}

type Address struct {
	Model
	UserId      uint
	HouseNo     uint   `gorm:"not null"`
	StreetName  string `gorm:"not null"`
	CityName    string `gorm:"not null"`
	StateName   string `gorm:"not null"`
	CountryName string `gorm:"not null"`
	ZipCode     string `gorm:"not null"`
	isArchived  bool   `gorm:"default:false"`
}

type Roles struct {
	Model
	UserId uint
	User   Users     `gorm:"foreignKey:UserId"`
	Role   RolesType `gorm:"type:role_type"`
}

type UserBody struct {
	Id       uint      `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	MobileNo string    `json:"mobileNo"`
	Role     RolesType `json:"role"`
}

type AddressBody struct {
	UserId      uint   `json:"userId"`
	HouseNo     uint   `json:"houseNo"`
	StreetName  string `json:"streetName"`
	CityName    string `json:"cityName"`
	StateName   string `json:"stateName"`
	CountryName string `json:"countryName"`
	ZipCode     string `json:"zipCode"`
}
