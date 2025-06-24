package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq" // Import the postgres driver

	"github.com/imhasandl/message-service/cmd/helper"
	"github.com/imhasandl/message-service/cmd/server"
	"github.com/imhasandl/message-service/internal/database"
	"github.com/imhasandl/message-service/internal/rabbitmq"
	"github.com/imhasandl/message-service/internal/redis"
	pb "github.com/imhasandl/message-service/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	env := helper.GetENVSecrets()

	lis, err := net.Listen("tcp", env.Port)
	if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}

	dbConn, err := sql.Open("postgres", env.DBURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)
	defer dbConn.Close()

	redisConfig := redis.NewRedisConfig(env.RedisSecret)
	redis.InitRedisClient(redisConfig)

	rabbitmq, err := rabbitmq.NewRabbitMQ(env.RabbitMQ)
	if err != nil {
		log.Fatalf("error initializing rabbit mq: %v", err)
	}
	defer rabbitmq.Close()

	server := server.NewServer(dbQueries, env.TokenSecret, rabbitmq)

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, server)

	reflection.Register(s)
	log.Printf("Server listening on %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to lister: %v", err)
	}
}
