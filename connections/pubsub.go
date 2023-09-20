package connections

import (
	"encoding/json"
	"fmt"

	"github.com/just-arun/micro-api-gateway/boot"
	"github.com/nats-io/nats.go"
)

func Pubsub(token string) {
	con := boot.NatsConnection(token)

	con.Subscribe("change-service-map", func(m *nats.Msg) {
		va := []boot.MapPathType{}
		err := json.Unmarshal(m.Data, &va)
		if err != nil {
			fmt.Println("ERR: ", err)
		}
		boot.MapPath = va
		m.Respond([]byte("answer is 42"))
	})
}
