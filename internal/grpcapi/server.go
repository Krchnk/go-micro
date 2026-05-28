package grpcapi

import (
	"context"
	"errors"
	log "log"
	"os"
	"strings"

	"github.com/Krchnk/go-micro/internal/notify"

	usersv1 "github.com/Krchnk/go-micro/internal/gen/users/v1"
	"github.com/Krchnk/go-micro/internal/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	usersv1.UnimplementedUserServiceServer
	service *users.Service
}

func NewServer(service *users.Service) *Server {
	return &Server{service: service}
}

func (s *Server) CreateUser(_ context.Context, req *usersv1.CreateUserRequest) (*usersv1.CreateUserResponse, error) {
	name := strings.TrimSpace(req.GetName())
	email := strings.TrimSpace(req.GetEmail())
	if name == "" || email == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email are required")
	}

	created := s.service.Create(name, email)

	// Публикация события в Kafka
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "user-registered"
	}
	err := notify.SendUserRegistered(broker, topic, notify.UserEvent{
		ID:    created.ID,
		Email: created.Email,
		Name:  created.Name,
	})
	if err != nil {
		log.Printf("Kafka publish error: %v", err)
	}
	return &usersv1.CreateUserResponse{User: toProtoUser(created)}, nil
}

func (s *Server) UpdateUser(_ context.Context, req *usersv1.UpdateUserRequest) (*usersv1.UpdateUserResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be positive")
	}

	name := strings.TrimSpace(req.GetName())
	email := strings.TrimSpace(req.GetEmail())
	if name == "" || email == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email are required")
	}

	updated, err := s.service.Update(req.GetId(), name, email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &usersv1.UpdateUserResponse{User: toProtoUser(updated)}, nil
}

func (s *Server) DeleteUser(_ context.Context, req *usersv1.DeleteUserRequest) (*usersv1.DeleteUserResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be positive")
	}

	if err := s.service.Delete(req.GetId()); err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &usersv1.DeleteUserResponse{}, nil
}

func (s *Server) ListUsers(_ context.Context, _ *usersv1.ListUsersRequest) (*usersv1.ListUsersResponse, error) {
	domainUsers := s.service.List()
	protoUsers := make([]*usersv1.User, 0, len(domainUsers))
	for _, user := range domainUsers {
		protoUsers = append(protoUsers, toProtoUser(user))
	}

	return &usersv1.ListUsersResponse{Users: protoUsers}, nil
}


// Auth реализует аутентификацию пользователя и возвращает JWT-токен
func (s *Server) Auth(_ context.Context, req *usersv1.AuthRequest) (*usersv1.AuthResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	password := strings.TrimSpace(req.GetPassword())
	if username == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password required")
	}

	// Примитивная проверка: любой пользователь с паролем "password" проходит
	if password != "password" {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	token, err := GenerateJWT(username)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &usersv1.AuthResponse{Token: token}, nil
}

func toProtoUser(user users.User) *usersv1.User {
	return &usersv1.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
