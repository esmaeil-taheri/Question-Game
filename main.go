package main

import (
	"time"

	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/validator/uservalidator"

	// "gameapp/repository/migrator"
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

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey: JwtSignKey,
			AccessExpirationTime: AccessTokenExpireDuration,
			RefreshExpirationTime: RefereshTokenExpireDuration,
			AccessSubject: AccessTokenSubject,
			RefreshSubject: RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "gameapp_user",
			Password: "gameapp_pass",
			Port: 3306,
			Host: "localhost",
			DBName: "gameapp",

		},
	}

	// TODO - add command for migrations
	// mgr := migrator.New(cfg.Mysql)
	// mgr.Up()

	authSvc, userSvc, userValidator := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(authSvc, MysqlRepo)

	uV := uservalidator.New(MysqlRepo)

	return authSvc, userSvc, uV
}
