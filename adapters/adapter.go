package adapters

import (
	"errors"

	"github.com/akshay0074700747/cart-service/entities"
	"gorm.io/gorm"
)

type CartAdapter struct {
	DB *gorm.DB
}

func NewCartAdapter(db *gorm.DB) *CartAdapter {
	return &CartAdapter{
		DB: db,
	}
}

func (cart *CartAdapter) CreateCart(req entities.Cart) (entities.Cart, error) {

	var res entities.Cart
	query := "INSERT INTO carts (user_id) VALUES($1) RETURNING id,user_id"

	tx := cart.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := cart.DB.Raw(query, req.UserID).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return entities.Cart{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (cart *CartAdapter) InsertIntoCart(req entities.CartItems, user_id uint) (entities.CartItems, error) {

	var res entities.CartItems
	query := "INSERT INTO cart_items (cart_id,product_id,quantity) SELECT c.id AS cart_id, $1 AS product_id, $2 AS quantity FROM carts c WHERE user_id = $3 RETURNING id,cart_id,product_id,quantity"

	tx := cart.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := cart.DB.Raw(query, req.ProductID, req.Quantity, user_id).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return entities.CartItems{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (cart *CartAdapter) GetCartItems(user_id uint) ([]entities.CartItems, error) {

	var res []entities.CartItems
	query := "SELECT * FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)"

	return res, cart.DB.Raw(query, user_id).Scan(&res).Error
}

func (cart *CartAdapter) DeleteCartItem(req entities.CartItems, user_id uint) error {

	query := "DELETE FROM cart_items WHERE product_id = $1 AND cart_id = (SELECT id FROM carts WHERE user_id = $2) "
	res := cart.DB.Exec(query, req.ProductID, user_id)

	tx := cart.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("Cart Itm not deleted")
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (cart *CartAdapter) TruncateCartItems(user_id uint) (error) {

	query := "DELETE FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1) RETURNING id,product_id,quantity"

	tx := cart.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := cart.DB.Raw(query, user_id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (cart *CartAdapter) IncrementQty(req entities.CartItems, user_id uint) (entities.CartItems, error) {

	var res entities.CartItems
	query := "UPDATE cart_items SET quantity = quantity + $1 WHERE product_id = $2 AND cart_id = (SELECT id FROM carts WHERE user_id = $3) RETURNS id,cart_id,product_id,quantity"

	tx := cart.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := cart.DB.Raw(query, req.Quantity, req.ProductID, user_id).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return entities.CartItems{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (cart *CartAdapter) DecrementQty(req entities.CartItems, user_id uint) (entities.CartItems, error) {

	var res entities.CartItems
	query := "UPDATE cart_items SET quantity = quantity - $1 WHERE product_id = $2 AND cart_id = (SELECT id FROM carts WHERE user_id = $3) RETURNING id,cart_id,product_id,quantity"

	tx := cart.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := cart.DB.Raw(query, req.Quantity, req.ProductID, user_id).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return entities.CartItems{}, err
	}
	if res.Quantity < 0 {
		tx.Rollback()
		return entities.CartItems{}, errors.New("the quantity cannot go below 0")
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}
