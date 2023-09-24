package main

import (
	"flag"

	"github.com/just-arun/micro-api-gateway/boot"
	"github.com/just-arun/micro-api-gateway/connections"
	"github.com/just-arun/micro-api-gateway/model"
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
	connections.Pubsub(env.Nats.Token)
	conn := boot.NewGrpcConnection(env.Grpc.Host, env.Grpc.Port)
	client := pb.NewSessionServiceClient(conn)
	server.Proxy(appPort, client)

}
