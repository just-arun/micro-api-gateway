package boot

type MapPathType struct {
	Key   string
	Value string
	Auth  bool
}

var MapPath = []MapPathType{
	{
		Key:   "auth",
		Value: "http://localhost:8090/api/v1",
		Auth:  false,
	},
}
