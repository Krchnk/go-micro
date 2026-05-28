package main

import (
	"context"
	"log"
	"os"
	"time"

	usersv1 "github.com/Krchnk/go-micro/internal/gen/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	addr := getEnv("GRPC_ADDR", "localhost:9090")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := usersv1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Получаем JWT токен через Auth
	authResp, err := client.Auth(ctx, &usersv1.AuthRequest{
		Username: "Ivan",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("Auth failed: %v", err)
	}
	token := authResp.GetToken()

	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	createdResp, err := client.CreateUser(ctx, &usersv1.CreateUserRequest{
		Name:  "Ivan",
		Email: "ivan@example.com",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	log.Printf("CreateUser: %+v", createdResp.GetUser())

	listBeforeUpdate, err := client.ListUsers(ctx, &usersv1.ListUsersRequest{})
	if err != nil {
		log.Fatalf("ListUsers (before update) failed: %v", err)
	}
	log.Printf("ListUsers before update: %+v", listBeforeUpdate.GetUsers())

	updatedResp, err := client.UpdateUser(ctx, &usersv1.UpdateUserRequest{
		Id:    createdResp.GetUser().GetId(),
		Name:  "Ivan Petrov",
		Email: "ivan.petrov@example.com",
	})
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	log.Printf("UpdateUser: %+v", updatedResp.GetUser())

	if _, err := client.DeleteUser(ctx, &usersv1.DeleteUserRequest{Id: createdResp.GetUser().GetId()}); err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	log.Printf("DeleteUser: id=%d", createdResp.GetUser().GetId())

	listAfterDelete, err := client.ListUsers(ctx, &usersv1.ListUsersRequest{})
	if err != nil {
		log.Fatalf("ListUsers (after delete) failed: %v", err)
	}
	log.Printf("ListUsers after delete: %+v", listAfterDelete.GetUsers())
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
