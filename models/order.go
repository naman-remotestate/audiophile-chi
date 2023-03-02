package models

type OrderStatus string

const (
	Completed    OrderStatus = "completed"
	Pending      OrderStatus = "pending"
	Cancelled    OrderStatus = "cancelled"
	NotDelivered OrderStatus = "not delivered"
)

type Orders struct {
	Model
	UserId      uint
	User        Users `gorm:"foreignKey:UserId"`
	AddressId   uint
	Address     Address     `gorm:"foreignKey:AddressId"`
	TotalCost   uint        `gorm:"not null"`
	OrderStatus OrderStatus `gorm:"type:order_status"`
}
