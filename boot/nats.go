package boot

import "github.com/nats-io/nats.go"

func NatsConnection(token string) *nats.EncodedConn {
	nc, err := nats.Connect(nats.DefaultURL, nats.Token(token))
	if err != nil {
		panic(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	return ec
}
