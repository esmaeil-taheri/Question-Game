package config

import (
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth authservice.Config
	Mysql mysql.Config
}
