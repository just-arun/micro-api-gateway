package main

import (
	"flag"

	"github.com/just-arun/micro-api-gateway/connections"
	"github.com/just-arun/micro-api-gateway/model"
	"github.com/just-arun/micro-api-gateway/server"
	"github.com/just-arun/micro-api-gateway/util"
)

var (
	appEnv  string
	appPort string
)

func init() {
	flag.StringVar(&appEnv, "env", "dev", "environment")
	flag.StringVar(&appPort, "port", ":3000", "environment")
	flag.Parse()
}

func main() {

	env := &model.Env{}
	util.GetEnv(".env."+appEnv, ".", &env)
	connections.Pubsub(env.Nats.Token)

	server.Proxy(appPort)
}
