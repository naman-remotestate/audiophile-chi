package models

import "time"

type Session struct {
	Model
	UserId    uint
	User      Users `gorm:"foreignKey:UserId"`
	StartTime time.Time
	EndTime   time.Time `gorm:"default:null"`
}

type SessionBody struct {
	Id        uint      `json:"id"`
	UserId    uint      `json:"userId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
