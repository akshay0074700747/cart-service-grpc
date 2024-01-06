package main

import (
	"log"
	"net"
	"os"

	"github.com/akshay0074700747/cart-service/db"
	"github.com/akshay0074700747/cart-service/initializer"
	"github.com/akshay0074700747/cart-service/service"
	"github.com/akshay0074700747/proto-files-for-microservices/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	if godotenv.Load(".env") != nil {
		log.Fatal("couldnt load the .env file")
	}

	dbAddr := os.Getenv("DATABASE_ADDR")
	if dbAddr == "" {
		log.Fatal("empty db address")
	}

	DB, err := db.InitDB(dbAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	servicee := initializer.InitAll(DB)
	server := grpc.NewServer()

	pb.RegisterCartServiceServer(server, servicee)

	productConn, err := grpc.Dial("product-service:50004", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}

	wishlistConn, err := grpc.Dial("wishlist-service:50007", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}

	service.InitClients(pb.NewProductServiceClient(productConn), pb.NewWishlistServiceClient(wishlistConn))

	listener, err := net.Listen("tcp", ":50006")

	if err != nil {
		log.Fatalf("Failed to listen on port 50002: %v", err)
	}

	log.Printf("Cart Server is listening on port")
	log.Println("i am running on k8s")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}
