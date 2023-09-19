package main

var MapPath = []struct {
	Key   string
	Value string
}{{
	Key:   "auth",
	Value: "http://localhost:8090/api/v1",
}}

func main() {
	Proxy(":8080")
}
