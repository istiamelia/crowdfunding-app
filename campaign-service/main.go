package main

import (
	"campaign-service/config"
	"campaign-service/gen/go/campaign/v1"
	"campaign-service/mq"
	"campaign-service/repository"
	"campaign-service/service"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main(){
	// Load .env only if running locally
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Warning: .env file not found, using environment variables instead")
		}
	}

	// Initialize database connection
	config.InitDB()
	gorm := config.DB

	//  Initialize rabbitMQ
	mq.InitRabbitMQ()
	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Create repository instances
	campaignRepo := repository.NewCampaignRepository(gorm)

	// Inject repositories into services
	campaignService := service.NewCampaignService(campaignRepo)

	// Register server with grpc
	campaign.RegisterCampaignServiceServer(grpcServer,campaignService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5051"
	}
	// Listener for grpcServer without interceptors
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on PORT: %v", err)
	}
	
	log.Printf("Starting gRPC server on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}