package server

import (
	"context"

	"github.com/imhasandl/message-service/internal/database"
	pb "github.com/imhasandl/message-service/protos"
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

func (s *server) SendMessage (ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return nil, nil
}

func (s *server) GetMessages (ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	return nil, nil
}