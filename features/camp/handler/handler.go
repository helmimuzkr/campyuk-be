package handler

import (
	"campyuk-api/features/camp"
	"campyuk-api/helper"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type campHandler struct {
	srv camp.CampService
}

func New(s camp.CampService) camp.CampHandler {
	return &campHandler{srv: s}
}

func (ch *campHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		cr := campRequest{}
		if err := c.Bind(&cr); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		form, err := c.MultipartForm()
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		documentHeader, found := form.File["document"]
		if !found {
			return c.JSON(helper.ErrorResponse("document input not found"))
		}
		imagesHeader, found := form.File["images"]
		if !found {
			return c.JSON(helper.ErrorResponse("images input not found"))
		}

		newCamp := camp.Core{}
		copier.Copy(&newCamp, &cr)

		if err := ch.srv.Add(token, newCamp, documentHeader[0], imagesHeader); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(201, "sukses add new camp"))
	}
}

func (ch *campHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		if token == nil {
			token = jwt.New(jwt.SigningMethodES256)
		}

		res, err := ch.srv.List(token)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := []listCampResponse{}
		copier.Copy(&response, &res)

		return c.JSON(helper.SuccessResponse(200, "success show list camp", response))
	}
}

func (ch *campHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
func (ch *campHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
func (ch *campHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
func (ch *campHandler) Accept() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
func (ch *campHandler) Decline() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
