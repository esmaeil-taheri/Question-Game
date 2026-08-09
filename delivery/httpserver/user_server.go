package httpserver

import (
	"gameapp/service/userservice"
	
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userRegister(c echo.Context) error {

	var Req userservice.RegisterRequest
	if err := c.Bind(&Req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	res, err := s.userSvc.Register(Req)
	if err != nil {
		return  echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}	

	return c.JSON(http.StatusCreated, res)
}

func (s Server) userLogin(c echo.Context) error {

	var Req userservice.LoginRequest
	if err := c.Bind(&Req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	res, err := s.userSvc.Login(Req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return  c.JSON(http.StatusOK, res)
}

func (s Server) userProfile(c echo.Context) error {

	authToken := c.Request().Header.Get("authorization")
	claims, err := s.authSvc.ParseToken(authToken)
	if err != nil {
		return  echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	res, err := s.userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		return  echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)

}
