package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Proxy(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.String(), "/")
		mapKey := ""
		mapValue := ""
		for _, v := range MapPath {
			if v.Key == path[1] {
				mapKey = v.Key
				mapValue = v.Value
				break
			}
		}
		urlValue := strings.Replace(r.URL.String(), "/"+mapKey, mapValue, 1)
		client := &http.Client{}
		if strings.Index(urlValue, "http") < 0 {
			w.WriteHeader(404)
			w.Write([]byte("404"))
			return
		}
		req, err := http.NewRequest("GET", urlValue, r.Body)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	http.ListenAndServe(port, nil)
}
