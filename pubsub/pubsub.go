package pubsub

import (
	"fmt"

	"github.com/just-arun/micro-api-gateway/boot"
	"github.com/just-arun/micro-api-gateway/model"
)

func Pubsub(token string) {
	con := boot.NatsConnection(token)

	con.Subscribe("change-service-map", func(m *[]model.ServiceMap) {
		fmt.Println(m)
		boot.MapPath = *m
		// va := []model.ServiceMap{}
		// fmt.Println(string(m.Data))
		// err := json.Unmarshal(m.Data, &va)
		// if err != nil {
		// 	fmt.Println("ERR: ", err)
		// 	return
		// }
		// boot.MapPath = va
		// fmt.Println(boot.MapPath)
		// m.Respond([]byte("ok"))
	})
}
