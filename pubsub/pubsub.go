package pubsub

import (
	"fmt"

	"github.com/just-arun/micro-api-gateway/boot"
	"github.com/just-arun/micro-api-gateway/model"
	"github.com/just-arun/micro-api-gateway/util"
	"github.com/nats-io/nats.go"
)

type pubsub struct {
	*nats.EncodedConn
}

func Pubsub(con *nats.EncodedConn) {
	ps := &pubsub{con}
	ps.Subscribe("change-service-map", func(m *[]model.ServiceMap) {
		if len(*m) < 1 {
			util.ReadJson(".", "sitemap.json", &boot.MapPath)
			return
		}
		boot.MapPath = *m
		for _, v := range boot.MapPath {
			fmt.Println(v)
		}
	})
}
