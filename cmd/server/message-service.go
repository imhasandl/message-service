package server

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/imhasandl/message-service/cmd/helper"
	"github.com/imhasandl/message-service/internal/database"
	"github.com/imhasandl/message-service/internal/rabbitmq"
	"github.com/imhasandl/message-service/internal/redis"
	pb "github.com/imhasandl/message-service/protos"
	"github.com/imhasandl/post-service/cmd/auth"
	postService "github.com/imhasandl/post-service/cmd/helper"
	"github.com/streadway/amqp"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server represents the gRPC server for the search service.
type Server interface {
	pb.MessageServiceServer
}

type server struct {
	pb.UnimplementedMessageServiceServer
	db          *database.Queries
	tokenSecret string
	rabbitmq    *rabbitmq.RabbitMQ
}

// NewServer creates and returns a new instance of the search service server.
// It requires database queries implementation and a token secret for authentication.
func NewServer(db *database.Queries, tokenSecret string, rabbitmq *rabbitmq.RabbitMQ) Server {
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
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't get token from header - SendMessage", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't get user id from token - SendMessage", err)
	}

	receiverID, err := uuid.Parse(req.GetReceiverId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse receiver's id to uuid - SendMessage", err)
	}

	var senderUserData database.User
	err = redis.GetCachedUser(userID.String(), &senderUserData)
	if err != nil {
		senderUserData, err = s.db.GetUserByID(ctx, userID)
		if err != nil {
			return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can;t get sender's data by id - SendMessage", err)
		}

		redis.CacheUser(userID.String(), senderUserData)
	}

	sendMessageParams := database.SendMessageParams{
		ID:         uuid.New(),
		SenderID:   userID,
		ReceiverID: receiverID,
		Content:    req.GetContent(),
	}

	message, err := s.db.SendMessage(ctx, sendMessageParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't send message via db - SendMessage", err)
	}

	redis.InvalidateMessagesCache(userID.String(), receiverID.String())
	redis.InvalidateConversationList(userID.String())
	redis.InvalidateConversationList(receiverID.String())
	redis.InvalidateLastMessage(userID.String(), receiverID.String())
	redis.DeleteMessageCount(userID.String(), receiverID.String())

	redis.CacheLastMessage(userID.String(), receiverID.String(), message)

	messageJSON, err := json.Marshal(map[string]interface{}{
		"title":           "New Notification",
		"sender_username": senderUserData.Username,
		"receiver_id":     receiverID.String(),
		"content":         req.GetContent(),
		"sent_at":         message.SentAt,
	})
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't marshal message to JSON - SendMessage", err)
	}

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
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't publish message to RabbitMQ - SendMessage", err)
	}

	return &pb.SendMessageResponse{
		Success: true,
	}, nil
}

func (s *server) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	accessToken, err := auth.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get user token from header - GetMessages", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't get user from given token - GetMessages", err)
	}

	receiverID, err := uuid.Parse(req.GetReceiverId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse receiver's id to uuid - GetMessages", err)
	}

	var messages []database.Message
	err = redis.GetCachedMessages(userID.String(), receiverID.String(), &messages)
	if err != nil {
		getMessagesParams := database.GetMessagesParams{
			SenderID:   userID,
			ReceiverID: receiverID,
		}

		messages, err = s.db.GetMessages(ctx, getMessagesParams)
		if err != nil {
			return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get messages from db - GetMessages", err)
		}

		redis.CacheMessages(userID.String(), receiverID.String(), messages)

		redis.CacheMessageCount(userID.String(), receiverID.String(), int64(len(messages)))
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

func (s *server) ChangeMessage(ctx context.Context, req *pb.ChangeMessageRequest) (*pb.ChangeMessageResponse, error) {
	messageID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse message id - ChangeMessage", err)
	}

	changeMessageParams := database.ChangeMessageParams{
		ID:      messageID,
		Content: req.GetContent(),
	}

	message, err := s.db.ChangeMessage(ctx, changeMessageParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't change message - ChangeMessage", err)
	}

	redis.InvalidateMessagesCache(message.SenderID.String(), message.ReceiverID.String())
	redis.InvalidateLastMessage(message.SenderID.String(), message.ReceiverID.String())

	return &pb.ChangeMessageResponse{
		Message: &pb.Message{
			Id:         message.ID.String(),
			SentAt:     timestamppb.New(message.SentAt),
			SenderId:   message.SenderID.String(),
			ReceiverId: message.ReceiverID.String(),
			Content:    message.Content,
		},
	}, nil
}

func (s *server) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	messageID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse message's id - DeleteMessage", err)
	}

	err = s.db.DeleteMessage(ctx, messageID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't delete message - DeleteMessage", err)
	}

	return &pb.DeleteMessageResponse{
		Status: true,
	}, nil
}
