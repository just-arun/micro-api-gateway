package grpcservice

import (
	"context"
	"fmt"

	"github.com/just-arun/micro-api-gateway/boot"
	"github.com/just-arun/micro-api-gateway/model"
	pb "github.com/just-arun/micro-session-proto"
)

type sitemap struct{}

func Sitemap() sitemap {
	return sitemap{}
}

func (st sitemap) GetServiceMap(client pb.SessionServiceClient) (err error) {
	stream, err := client.GetServiceMap(context.Background(), &pb.NoPayload{})
	if err != nil {
		fmt.Println("ERR: ", err.Error())
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			return nil
		}
		ca := &model.ServiceMap{
			ID:    uint(resp.Id),
			Key:   resp.Key,
			Value: resp.Value,
			Auth:  resp.Auth,
		}
		boot.MapPath = append(boot.MapPath, *ca)
	}
}
