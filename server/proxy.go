package server

import (
	"io"
	"net/http"
	"strings"

	"github.com/just-arun/micro-api-gateway/boot"
)

func getSortedUrl(r *http.Request) string {
	path := strings.Split(r.URL.String(), "/")
	mapKey := ""
	mapValue := ""
	for _, v := range boot.MapPath {
		if v.Key == path[1] {
			mapKey = v.Key
			mapValue = v.Value
			break
		}
	}
	return strings.Replace(r.URL.String(), "/"+mapKey, mapValue, 1)
}

func copyRequestHeader(req *http.Request, r *http.Request) {
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
}

func copyResponseHeader(w http.ResponseWriter, resp *http.Response) {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
}

func Proxy(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlValue := getSortedUrl(r)
		client := &http.Client{}
		containInMap := strings.Index(urlValue, "http")
		if containInMap < 0 {
			w.WriteHeader(404)
			w.Write([]byte("404"))
			return
		}
		req, err := http.NewRequest(r.Method, urlValue, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		copyRequestHeader(req, r)

		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		copyResponseHeader(w, resp)
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	http.ListenAndServe(port, nil)
}
