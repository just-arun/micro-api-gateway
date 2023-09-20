package boot

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcConnection(host, port string) *grpc.ClientConn {
	conn, err := grpc.Dial(host+port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return conn
}
