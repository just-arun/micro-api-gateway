package model

type grpc struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type natsEnv struct {
	Token string `mapstructure:"token"`
}

type Env struct {
	Grpc grpc    `mapstructure:"grpc"`
	Nats natsEnv `mapstructure:"nats"`
}
