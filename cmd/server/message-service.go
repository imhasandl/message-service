package server

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/imhasandl/message-service/cmd/helper"
	"github.com/imhasandl/message-service/internal/database"
	"github.com/imhasandl/message-service/internal/rabbitmq"
	pb "github.com/imhasandl/message-service/protos"
	"github.com/imhasandl/post-service/cmd/auth"
	postService "github.com/imhasandl/post-service/cmd/helper"
	"github.com/streadway/amqp"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedMessageServiceServer
	db          *database.Queries
	tokenSecret string
	rabbitmq    *rabbitmq.RabbitMQ
}

func NewServer(db *database.Queries, tokenSecret string, rabbitmq *rabbitmq.RabbitMQ) *server {
	return &server{
		pb.UnimplementedMessageServiceServer{},
		db,
		tokenSecret,
		rabbitmq,
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
		ID:         uuid.New(),
		SenderID:   userID,
		ReceiverID: receiverID,
		Content:    req.GetContent(),
	}

	message, err := s.db.SendMessage(ctx, sendMessageParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't send message via db", err)
	}

	// Create JSON message
	messageJSON, err := json.Marshal(map[string]interface{}{
		"receiver_id": receiverID.String(),
		"content":     req.GetContent(),
		"sent_at":     message.SentAt,
	})
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't marshal message to JSON", err)
	}

	// Publish message to RabbitMQ
	err = s.rabbitmq.Channel.Publish(
		rabbitmq.ExchangeName, // exchange
		rabbitmq.RoutingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		})
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't publish message to RabbitMQ", err)
	}

	return &pb.SendMessageResponse{}, nil
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
		SenderID:   userID,
		ReceiverID: receiverID,
	}

	messages, err := s.db.GetMessages(ctx, getMessagesParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't  get messages from db", err)
	}

	messagesResponse := make([]*pb.Message, len(messages))
	for i, message := range messages {
		messagesResponse[i] = &pb.Message{
			Id:         message.ID.String(),
			SentAt:     timestamppb.New(message.SentAt),
			SenderId:   message.SenderID.String(),
			ReceiverId: message.ReceiverID.String(),
			Content:    message.Content,
		}
	}

	return &pb.GetMessagesResponse{
		Message: messagesResponse,
	}, nil
}
