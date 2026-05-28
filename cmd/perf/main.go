package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	usersv1 "github.com/Krchnk/go-micro/internal/gen/users/v1"
	"google.golang.org/protobuf/proto"
)

type jsonUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type jsonPayload struct {
	Users []jsonUser `json:"users"`
}

func main() {
	iterations := 100000
	protoData := &usersv1.ListUsersResponse{Users: make([]*usersv1.User, 0, 100)}
	jsonData := jsonPayload{Users: make([]jsonUser, 0, 100)}

	for i := range 100 {
		id := int64(i + 1)
		name := fmt.Sprintf("User-%d", i+1)
		email := fmt.Sprintf("user-%d@example.com", i+1)

		protoData.Users = append(protoData.Users, &usersv1.User{Id: id, Name: name, Email: email})
		jsonData.Users = append(jsonData.Users, jsonUser{ID: id, Name: name, Email: email})
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("json marshal failed: %v", err)
	}

	protoBytes, err := proto.Marshal(protoData)
	if err != nil {
		log.Fatalf("proto marshal failed: %v", err)
	}

	startJSONMarshal := time.Now()
	for i := 0; i < iterations; i++ {
		if _, err := json.Marshal(jsonData); err != nil {
			log.Fatalf("json marshal loop failed: %v", err)
		}
	}
	jsonMarshalDuration := time.Since(startJSONMarshal)

	startProtoMarshal := time.Now()
	for i := 0; i < iterations; i++ {
		if _, err := proto.Marshal(protoData); err != nil {
			log.Fatalf("proto marshal loop failed: %v", err)
		}
	}
	protoMarshalDuration := time.Since(startProtoMarshal)

	startJSONUnmarshal := time.Now()
	for i := 0; i < iterations; i++ {
		var out jsonPayload
		if err := json.Unmarshal(jsonBytes, &out); err != nil {
			log.Fatalf("json unmarshal loop failed: %v", err)
		}
	}
	jsonUnmarshalDuration := time.Since(startJSONUnmarshal)

	startProtoUnmarshal := time.Now()
	for i := 0; i < iterations; i++ {
		var out usersv1.ListUsersResponse
		if err := proto.Unmarshal(protoBytes, &out); err != nil {
			log.Fatalf("proto unmarshal loop failed: %v", err)
		}
	}
	protoUnmarshalDuration := time.Since(startProtoUnmarshal)

	fmt.Printf("Payload size (JSON): %d bytes\n", len(jsonBytes))
	fmt.Printf("Payload size (Proto): %d bytes\n", len(protoBytes))
	fmt.Printf("Marshal JSON: %v\n", jsonMarshalDuration)
	fmt.Printf("Marshal Proto: %v\n", protoMarshalDuration)
	fmt.Printf("Unmarshal JSON: %v\n", jsonUnmarshalDuration)
	fmt.Printf("Unmarshal Proto: %v\n", protoUnmarshalDuration)
}
