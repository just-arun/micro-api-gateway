package grpcservice

import (
	"context"
	"fmt"
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
	fmt.Println("TOKEN", accessToken)
	result, err := client.VerifyUserSession(ctx, &pb.VerifyUserSessionParams{
		Token: accessToken,
	})
	if err != nil {
	fmt.Println("EE", err.Error())
		return nil, err
	}
	return result, nil
}

func (st session) GetUserSessionRefreshToken(client pb.SessionServiceClient, refreshToken string) (*pb.SetUserSessionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := client.GetUserSessionRefreshToken(ctx, &pb.GetUserSessionRefreshTokenPayload{
		RefreshToken: refreshToken,
		AccessTokenExpireInMinutes: 600,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}



