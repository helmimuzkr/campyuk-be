package handler

import (
	"campyuk-api/features/user"
	"campyuk-api/pkg/helper"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type userHandler struct {
	srv  user.UserService
	conf *oauth2.Config
}

func New(srv user.UserService, conf *oauth2.Config) user.UserHandler {
	return &userHandler{
		srv:  srv,
		conf: conf,
	}
}

func (uc *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := LoginRequest{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "input format incorrect")
		}
		if input.Username == "" {
			return c.JSON(helper.ErrorResponse("username is empty"))
		} else if input.Password == "" {
			return c.JSON(helper.ErrorResponse("password is empty"))
		}

		token, res, err := uc.srv.Login(input.Username, input.Password)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(http.StatusOK, "success login", ToResponse(res), token))
	}
}

func (uc *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := RegisterRequest{}
		err := c.Bind(&input)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "input format incorrect")
		}

		_, err = uc.srv.Register(*ReqToCore(input))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusCreated, "success create account"))
	}
}

func (uc *userHandler) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		res, err := uc.srv.Profile(token)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(http.StatusOK, "success show profile", GetToResponse(res)))
	}
}

func (uc *userHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := UpdateRequest{}
		err := c.Bind(&input)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "input format incorrect"})
		}

		formHeader, err := c.FormFile("user_image")
		if err != nil {
			log.Println(err)
		}

		_, err = uc.srv.Update(c.Get("user"), formHeader, *ReqToCore(input))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(http.StatusOK, "success update profile"))
	}
}

func (uc *userHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := uc.srv.Delete(c.Get("user"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "internal server error",
			})
		}
		return c.NoContent(204)
	}
}

func (uc *userHandler) GoogleAuth() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, uc.conf.AuthCodeURL("random", oauth2.AccessTypeOffline))
	}
}

func (uc *userHandler) GoogleCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")

		token, err := uc.conf.Exchange(c.Request().Context(), code)
		if err != nil {
			log.Println(err)
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		_, err = uc.srv.LoginGoogle(token.AccessToken, token.RefreshToken)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "login success"))
	}
}
