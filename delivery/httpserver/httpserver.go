package httpserver

import (
	"fmt"
	"log/slog"

	"gameapp/config"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
    config config.Config
    authSvc authservice.Service
    userSvc userservice.Service
    userValidator uservalidator.Validator
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Server {
    return Server{
        config: config,
        authSvc: authSvc,
        userSvc: userSvc,
        userValidator: userValidator,
    }
}

func (s Server) Serve() {
    e := echo.New()

   // Middleware
    e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
    e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)

    // Routes
    userGroup := e.Group("/users")

    userGroup.GET("/profile", s.userProfile)
    userGroup.POST("/login", s.userLogin)
    userGroup.POST("/register", s.userRegister)

    if err := e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)); err != nil {
        slog.Error("failed to start server", "error", err)
    }
}