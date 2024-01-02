package initializer

import (
	"github.com/akshay0074700747/cart-service/adapters"
	"github.com/akshay0074700747/cart-service/service"
	"gorm.io/gorm"
)

func InitAll(db *gorm.DB) *service.CartService {
	return service.NewCartService(adapters.NewCartAdapter(db))
}
