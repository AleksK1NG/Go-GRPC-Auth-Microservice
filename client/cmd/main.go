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
		"session_id", "e0dad68c-f612-42e5-85df-cd7e989330fd",
		"subsystem", "cli",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := client.GetMe(ctx, &userService.GetMeRequest{})
	if err != nil {
		log.Fatal(err)
	}

	//res, err := client.Logout(ctx, &userService.LogoutRequest{})
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Printf("RESPONSE: %s", res.String())
}
