package main

import (
	"context"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	grcpConn, err := grpc.Dial(
		"127.0.0.1:5000",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	client := userService.NewUserServiceClient(grcpConn)

	ctx := context.Background()
	md := metadata.Pairs(
		"session_id", "63dc333d-0863-4e9d-899e-bfc7a2ce9217",
		"subsystem", "cli",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := client.GetMe(ctx, &userService.GetMeRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("RESPONSE: %s", res.String())
}
