package entities

type Cart struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"unique"`
}

type CartItems struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint `gorm:"foreignKey:CartID;references:carts(id)"`
	ProductID uint
	Quantity  int
}
