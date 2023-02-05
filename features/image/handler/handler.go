package handler

import (
	"campyuk-api/features/image"
	"campyuk-api/helper"
	"strconv"

	"github.com/labstack/echo/v4"
)

type imageHandler struct {
	srv image.ImageService
}

func New(srv image.ImageService) image.ImageHandler {
	return &imageHandler{srv: srv}
}

func (ih *imageHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		header, err := c.FormFile("image")
		if err != nil {
			return c.JSON(helper.ErrorResponse("image input not found"))
		}

		str := c.FormValue("camp_id")
		campID, _ := strconv.Atoi(str)

		if err := ih.srv.Add(token, uint(campID), header); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success add new image"))
	}
}
func (ih *imageHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
func (ih *imageHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.Param("id")
		imageID, _ := strconv.Atoi(str)

		if err := ih.srv.Delete(token, uint(imageID)); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success delete image"))
	}
}
