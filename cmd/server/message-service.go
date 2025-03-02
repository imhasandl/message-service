package server

import (
	pb "github.com/imhasandl/message-service/protos"
	"github.com/imhasandl/message-service/internal/database"
)

type server struct {
	pb.UnimplementedMessageServiceServer
	db          *database.Queries
	tokenSecret string
}

func NewServer(db *database.Queries, tokenSecret string) *server {
	return &server{
		pb.UnimplementedMessageServiceServer{},
		db,
		tokenSecret,
	}
}