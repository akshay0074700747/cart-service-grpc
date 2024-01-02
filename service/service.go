package service

import (
	"context"
	"errors"
	"log"

	"github.com/akshay0074700747/cart-service/adapters"
	"github.com/akshay0074700747/cart-service/entities"
	"github.com/akshay0074700747/proto-files-for-microservices/pb"
)

var (
	ProductClient  pb.ProductServiceClient
	WishlistClient pb.WishlistServiceClient
)

func InitClients(product pb.ProductServiceClient, wish pb.WishlistServiceClient) {
	ProductClient = product
	WishlistClient = wish
}

type CartService struct {
	Adapter adapters.AdapterInterface
	pb.UnimplementedCartServiceServer
}

func NewCartService(adapter adapters.AdapterInterface) *CartService {
	return &CartService{
		Adapter: adapter,
	}
}

func (cart *CartService) CreateCart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponce, error) {

	res, err := cart.Adapter.CreateCart(entities.Cart{UserID: uint(req.UserId)})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.CartResponce{CartId: uint32(res.ID), UserId: uint32(res.UserID), IsEmpty: true}, nil
}

func (cart *CartService) GetCart(ctx context.Context, req *pb.CartRequest) (*pb.GetCartResponce, error) {

	var ids []uint32
	var cartProds []*pb.AddtoCartResponce

	res, err := cart.Adapter.GetCartItems(uint(req.UserId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if len(res) == 0 {
		return &pb.GetCartResponce{}, nil
	}

	for _, prod := range res {
		ids = append(ids, uint32(prod.ProductID))
	}

	productRes, err := ProductClient.GetArrayofProducts(ctx, &pb.ArrayofProductsRequest{Id: ids})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for i, prod := range productRes.Products {
		cartProds = append(cartProds, &pb.AddtoCartResponce{
			Quantity: int32(res[i].Quantity),
			Product:  prod,
		})
	}

	return &pb.GetCartResponce{CartId: uint32(res[0].CartID), Products: cartProds}, nil
}

func (cart *CartService) AddtoCart(ctx context.Context, req *pb.AddtoCartRequest) (*pb.AddtoCartResponce, error) {

	productRes, err := ProductClient.GetProduct(context.TODO(), &pb.GetProductByID{Id: req.ProductId})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New(err.Error())
	}

	if productRes.Name == "" {
		return nil, errors.New("Sorry the product doesnt exist")
	}

	item, err := cart.Adapter.InsertIntoCart(entities.CartItems{
		ProductID: uint(req.ProductId),
		Quantity:  int(req.Quantity)}, uint(req.UserId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.AddtoCartResponce{
		Product: &pb.AddProductResponce{
			Id:    productRes.Id,
			Name:  productRes.Name,
			Price: productRes.Price,
			Stock: productRes.Stock,
		},
		Quantity: int32(item.Quantity),
	}, nil
}

func (cart *CartService) DeleteCartItem(ctx context.Context, req *pb.AddtoCartRequest) (*pb.GetCartResponce, error) {

	productRes, err := ProductClient.GetProduct(context.TODO(), &pb.GetProductByID{Id: req.ProductId})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New(err.Error())
	}

	if productRes.Name == "" {
		return nil, errors.New("A product with the given ID doesnt exist")
	}

	err = cart.Adapter.DeleteCartItem(entities.CartItems{ProductID: uint(req.ProductId)}, uint(req.UserId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.GetCartResponce{}, nil
}

func (cart *CartService) TruncateCart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponce, error) {

	if err := cart.Adapter.TruncateCartItems(uint(req.UserId)); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.CartResponce{UserId: req.UserId, IsEmpty: true}, nil
}

func (cart *CartService) ChangeQty(ctx context.Context, req *pb.ChangeQtyRequest) (*pb.AddtoCartResponce, error) {

	var res entities.CartItems
	var err error
	if req.IsIncrease {

		res, err = cart.Adapter.IncrementQty(entities.CartItems{
			ProductID: uint(req.ProductId),
			Quantity:  int(req.Quantity),
		}, uint(req.UserId))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	} else {

		res, err = cart.Adapter.DecrementQty(entities.CartItems{
			ProductID: uint(req.ProductId),
			Quantity:  int(req.Quantity),
		}, uint(req.UserId))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	prod, err := ProductClient.GetProduct(ctx, &pb.GetProductByID{Id: uint32(res.ProductID)})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.AddtoCartResponce{Product: prod, Quantity: int32(res.Quantity)}, nil
}

func (cart *CartService) TrasferWishlist(ctx context.Context, req *pb.CartRequest) (*pb.GetCartResponce, error) {

	var result []*pb.AddtoCartResponce
	wishRes, err := WishlistClient.GetWishlist(ctx, &pb.WishlistRequest{UserId: req.UserId})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var cartid uint
	for _, wish := range wishRes.Products {
		cart, err := cart.Adapter.InsertIntoCart(entities.CartItems{
			ProductID: uint(wish.Id),
			Quantity:  1,
		}, uint(req.UserId))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		cartid = cart.ID
		result = append(result, &pb.AddtoCartResponce{Product: wish, Quantity: int32(cart.Quantity)})
	}

	return &pb.GetCartResponce{CartId: uint32(cartid), Products: result}, nil
}
