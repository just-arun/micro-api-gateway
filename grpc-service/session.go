package grpcservice

import (
	"context"
	"time"

	pb "github.com/just-arun/micro-session-proto"
)

type session struct{}

func Session() session {
	return session{}
}

func (st session) VerifySession(client pb.SessionServiceClient, accessToken string) (*pb.VerifyUserSessionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := client.VerifyUserSession(ctx, &pb.VerifyUserSessionParams{
		Token: accessToken,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}



