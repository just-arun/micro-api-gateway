package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/just-arun/micro-api-gateway/boot"
	grpcservice "github.com/just-arun/micro-api-gateway/grpc-service"
	"github.com/just-arun/micro-api-gateway/model"
	pb "github.com/just-arun/micro-session-proto"
)

func cors(r *http.Request, w http.ResponseWriter, env *model.Env) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	// util.Array().Includes(env.Cors.Origins, func(item string, index int) bool {return item == })
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
}

func getSortedData(r *http.Request) (data *model.ServiceMap, url string) {
	path := strings.Split(r.URL.String(), "/")
	mapKey := ""
	mapValue := ""
	fmt.Println(boot.MapPath)
	for _, v := range boot.MapPath {
		if v.Key == path[1] {
			mapKey = v.Key
			mapValue = v.Value
			data = &v
			break
		}
	}
	url = strings.Replace(r.URL.String(), "/"+mapKey, mapValue, 1)
	return
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

func authVerify(r *http.Request, conn pb.SessionServiceClient, req *http.Request) error {
	accessToken, err := r.Cookie("x-access")
	if err != nil {
		return err
	}
	verifyAccess, err := grpcservice.Session().VerifySession(conn, accessToken.Value)
	if err != nil {
		return err
	}

	req.Header.Del("x-roles")
	req.Header.Del("x-token")
	req.Header.Del("x-user-id")

	req.Header.Add("x-roles", strings.Join(verifyAccess.Roles, ","))
	req.Header.Add("x-token", accessToken.Value)
	req.Header.Add("x-user-id", strconv.FormatInt(int64(verifyAccess.UserID), 10))
	return nil
}

func Proxy(port string, conn pb.SessionServiceClient, env *model.Env) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cors(r, w, env)
		urlData, urlValue := getSortedData(r)
		fmt.Println(urlValue)
		client := &http.Client{}
		req, err := http.NewRequest(r.Method, urlValue, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		copyRequestHeader(req, r)

		// auth verify
		if urlData.Auth {
			err = authVerify(r, conn, req)
			if err != nil {
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized 401(0)"))
				return
			}
		}

		containInMap := strings.Index(urlValue, "http")
		if containInMap < 0 {
			w.WriteHeader(404)
			w.Write([]byte("404"))
			return
		}

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

	fmt.Println("serving proxy")
	http.ListenAndServe(port, nil)
}
