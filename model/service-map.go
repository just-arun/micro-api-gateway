package model

type ServiceMap struct {
	ID    uint   `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Auth  bool   `json:"auth"`
}