package main

import (
	"flag"

	"github.com/just-arun/micro-api-gateway/boot"
	grpcservice "github.com/just-arun/micro-api-gateway/grpc-service"
	"github.com/just-arun/micro-api-gateway/model"
	"github.com/just-arun/micro-api-gateway/pubsub"
	"github.com/just-arun/micro-api-gateway/server"
	"github.com/just-arun/micro-api-gateway/util"
	pb "github.com/just-arun/micro-session-proto"
)

var (
	appEnv  string
	appPort string
)

func init() {
	flag.StringVar(&appEnv, "env", "dev", "environment")
	flag.StringVar(&appPort, "port", ":8080", "environment")
	flag.Parse()
}

func main() {

	env := &model.Env{}
	util.GetEnv(".env."+appEnv, ".", &env)
	pubsub.Pubsub(env.Nats.Token)
	conn := boot.NewGrpcConnection(env.Grpc.Host, env.Grpc.Port)
	client := pb.NewSessionServiceClient(conn)
	err := grpcservice.Sitemap().GetServiceMap(client)
	if err != nil {
		panic(err)
	}
	server.Proxy(appPort, client)

}
