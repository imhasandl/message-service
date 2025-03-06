package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/imhasandl/message-service/cmd/helper"
	"github.com/imhasandl/message-service/internal/database"
	pb "github.com/imhasandl/message-service/protos"
	"github.com/imhasandl/post-service/cmd/auth"
	postService "github.com/imhasandl/post-service/cmd/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	accessToken, err := auth.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't get token from header", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't get user id from token", err)
	}

	receiverID, err := uuid.Parse(req.GetReceiverId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse receiver's id to uuid", err)
	}

	sendMessageParams := database.SendMessageParams{
		SenderID:   userID,
		ReceiverID: receiverID,
		Content:    req.GetContent(),
	}

	message, err := s.db.SendMessage(ctx, sendMessageParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't send message via db", err)
	}

	return &pb.SendMessageResponse{
		Message: &pb.Message{
			Id:         message.ID.String(),
			SentAt:     timestamppb.New(message.SentAt),
			SenderId:   message.SenderID.String(),
			ReceiverId: message.ReceiverID.String(),
			Content:    message.Content,
		},
	}, nil
}

func (s *server) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	accessToken, err := auth.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get user token from header", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't get user from given token", err)
	}

	receiverID, err := uuid.Parse(req.GetReceiverId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse receiver's id to uuid", err)
	}

	getMessagesParams := database.GetMessagesParams{
		SenderID: userID,
		ReceiverID: receiverID,
	}

	messages, err := s.db.GetMessages(ctx, getMessagesParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't  get messages from db", err)
	}

	messagesResponse := make([]*pb.Message, len(messages))
	for i, message := range messages {
		messagesResponse[i] = &pb.Message{
			Id: message.ID.String(),
			SentAt: timestamppb.New(message.SentAt),
			SenderId: message.SenderID.String(),
			ReceiverId: message.ReceiverID.String(),
			Content: message.Content,
		}
	}

	return &pb.GetMessagesResponse{
		Message: messagesResponse,
	}, nil
}
