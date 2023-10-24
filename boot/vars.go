package boot

import "github.com/just-arun/micro-api-gateway/model"

type MapPathType struct {
	Key   string
	Value string
	Auth  bool
}

var MapPath = []model.ServiceMap{}

var GeneralSettings = model.General{
	TokenPlacement: model.TokenPlacementHeader,
}
