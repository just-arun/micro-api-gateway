package main

import (
	"flag"
	"fmt"

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
	con := boot.NatsConnection(env.Nats.Token)
	pubsub.Pubsub(con)
	conn := boot.NewGrpcConnection(env.Grpc.Host, env.Grpc.Port)
	client := pb.NewSessionServiceClient(conn)
	err := grpcservice.Sitemap().GetServiceMap(client)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Trying backup")
		err = util.ReadJson(".", "sitemap.json", &boot.MapPath)
		if err != nil {
			panic("no backup sitemap file found")
		}
	}

	server.Proxy(appPort, client, env)

}
