package adapters

import "github.com/akshay0074700747/cart-service/entities"

type AdapterInterface interface {
	CreateCart(req entities.Cart) (entities.Cart, error)
	InsertIntoCart(req entities.CartItems, user_id uint) (entities.CartItems, error)
	GetCartItems(user_id uint) ([]entities.CartItems, error)
	DeleteCartItem(req entities.CartItems, user_id uint) error
	TruncateCartItems(user_id uint) error
	IncrementQty(req entities.CartItems, user_id uint) (entities.CartItems, error)
	DecrementQty(req entities.CartItems, user_id uint) (entities.CartItems, error)
}
