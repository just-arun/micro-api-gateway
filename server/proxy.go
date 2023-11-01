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
	"github.com/just-arun/micro-api-gateway/util"
	pb "github.com/just-arun/micro-session-proto"
)

func cors(r *http.Request, w http.ResponseWriter, env *model.Env) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token, x-refresh"
	includes := util.Array().Includes(env.Cors.Origins, func(item string, index int) bool { return item == r.Header.Get("Origin") })
	if includes {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
}

func getSortedData(r *http.Request) (data *model.ServiceMap, url string) {
	path := strings.Split(r.URL.String(), "/")
	mapKey := ""
	mapValue := ""
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

type foundHeaderTokenType int

const (
	foundHeaderTokenTypeAccessToken foundHeaderTokenType = iota
	foundHeaderTokenTypeRefreshToken
	foundHeaderTokenTypeNoneToken
)

func tokenFromHeader(r *http.Request) (access string, refresh string, tokenType foundHeaderTokenType, err error) {
	tokenType = foundHeaderTokenTypeNoneToken
	accessToken := r.Header.Get("Authorization")
	if len(accessToken) == 0 {
		refreshToken := r.Header.Get("x-refresh")
		if len(refreshToken) != 0 {
			return "", "", foundHeaderTokenTypeNoneToken, nil
		}
		return "", refreshToken, foundHeaderTokenTypeRefreshToken, nil
	}
	accessT := strings.Split(accessToken, " ")
	if accessT[0] != "barear" {
		return "", "", foundHeaderTokenTypeRefreshToken, fmt.Errorf("invalid authorization token")
	}
	return accessT[1], "", foundHeaderTokenTypeAccessToken, nil
}

func tokenFromCookie(r *http.Request) (access string, refresh string, tokenType foundHeaderTokenType) {
	tokenType = foundHeaderTokenTypeNoneToken
	accessToken, err := r.Cookie("x-session")
	if err != nil {
		refreshToken, err := r.Cookie("x-refresh")
		if err != nil {
			return "", "", foundHeaderTokenTypeNoneToken
		}
		return "", refreshToken.Value, foundHeaderTokenTypeRefreshToken
	}
	return accessToken.Value, "", foundHeaderTokenTypeAccessToken
}

func fetchAndSeedUserData(conn pb.SessionServiceClient, req *http.Request, accessToken string) (err error) {
	verifyAccess, err := grpcservice.Session().VerifySession(conn, accessToken)
	if err != nil {
		fmt.Println("DA", verifyAccess)
		return err
	}

	req.Header.Del("x-roles")
	req.Header.Del("x-token")
	req.Header.Del("x-user-id")

	req.Header.Add("x-roles", strings.Join(verifyAccess.Roles, ","))
	req.Header.Add("x-token", accessToken)
	req.Header.Add("x-user-id", strconv.FormatInt(int64(verifyAccess.UserID), 10))
	return
}

func fetchAndNewAccessToken(conn pb.SessionServiceClient, refreshToken string) (accessToken string, err error) {
	response, err := grpcservice.Session().GetUserSessionRefreshToken(conn, refreshToken)

	if err != nil {
		return accessToken, err
	}

	return response.Token, nil
}

func authVerify(r *http.Request, conn pb.SessionServiceClient, req *http.Request, w http.ResponseWriter) (err error) {
	accessToken := ""
	refreshToken := ""
	tokenType := foundHeaderTokenTypeNoneToken

	fmt.Println("T Type", tokenType)

	if boot.GeneralSettings.TokenPlacement == model.TokenPlacementHeader {
		accessToken, refreshToken, tokenType, err = tokenFromHeader(r)
		if err != nil {
			return fmt.Errorf("invalid token {0}")
		}
	}

	if boot.GeneralSettings.TokenPlacement == model.TokenPlacementCookie {
		accessToken, refreshToken, tokenType = tokenFromCookie(r)
	}
	
	fmt.Println("T Type", tokenType)

	if tokenType == foundHeaderTokenTypeNoneToken {
		return fmt.Errorf("no access token found")
	}

	switch tokenType {
	case foundHeaderTokenTypeAccessToken:
		fmt.Println("ca","da")
		return fetchAndSeedUserData(conn, req, accessToken)
	case foundHeaderTokenTypeRefreshToken:
		fmt.Println("12")
		accessToken, err = fetchAndNewAccessToken(conn, refreshToken)
		if err != nil {
		fmt.Println("EEx", err.Error())
			return err
		}
		fmt.Println("1ac", accessToken)
		if boot.GeneralSettings.TokenPlacement == model.TokenPlacementCookie {
			http.SetCookie(w, &http.Cookie{
				Name:   "x-session",
				Value:  accessToken,
				Path:   "/",
				Secure: true,
				MaxAge: int(boot.GeneralSettings.AccessTokenExpiryTime),
			})
		}
		return fetchAndSeedUserData(conn, req, accessToken)
	default:
		return fmt.Errorf("invalid token {112}")
	}
}

func Proxy(port string, conn pb.SessionServiceClient, env *model.Env) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cors(r, w, env)
		urlData, urlValue := getSortedData(r)
		client := &http.Client{}
		req, err := http.NewRequest(r.Method, urlValue, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		copyRequestHeader(req, r)

		fmt.Println(-1, "pre-auth", urlValue, urlData.Auth)
		// auth verify
		if urlData.Auth {
			fmt.Println(0)
			err = authVerify(r, conn, req, w)
			if err != nil {
				fmt.Println(1)
				w.WriteHeader(401)
				w.Write([]byte(`Unauthorized 401(0) ` + err.Error()))
				return
			}
		}

		containInMap := strings.Index(urlValue, "http")
		if containInMap < 0 {
			w.WriteHeader(404)
			w.Write([]byte("page not found"))
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
