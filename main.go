package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
)

const (
	JwtSignKey = "jwt_secret"
	AccessTokenSubject = "ac"
	RefreshTokenSubject = "rt"
	AccessTokenExpireDuration = time.Hour * 24
	RefereshTokenExpireDuration = time.Hour * 7 * 24
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healtCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)

	log.Println("Server is listening on port 8080...")
	server := http.Server{Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error": "invalid method"}`)

		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpireDuration, RefereshTokenExpireDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	_, err = userSvc.Register(uReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}	

	writer.Write([]byte(`{"message": "user created"}`))

}

func healtCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "healthy"}`)
}

func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error": "invalid method"}`)

		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	var lReq userservice.LoginRequest
	err = json.Unmarshal(data, &lReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpireDuration, RefereshTokenExpireDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	res, err := userSvc.Login(lReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	data, err = json.Marshal(res)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	writer.Write(data)

}

func userProfileHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodGet {
		fmt.Fprintf(writer, `{"error": "invalid method"}`)

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpireDuration, RefereshTokenExpireDuration)
	
	authToken := req.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(authToken)
	if err != nil {
		fmt.Fprintf(writer, `{"error": "token is not valid"}`)
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	res, err := userSvc.Profile(userservice.ProfileRequest{claims.UserID})
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	writer.Write(data)
}
