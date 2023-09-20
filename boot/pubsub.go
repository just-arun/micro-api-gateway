package boot

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

func Pubsub(token string) {
	con := NatsConnection(token)

	con.Subscribe("change-service-map", func(m *nats.Msg) {
		type d struct {
			Key   string
			Value string
		}
		va := []d{}
		err := json.Unmarshal(m.Data, &va)
		if err != nil {
			fmt.Println("ERR: ", err)
		}
		m.Respond([]byte("answer is 42"))
	})
}
