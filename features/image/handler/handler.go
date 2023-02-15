package handler

import (
	"campyuk-api/features/image"
	"campyuk-api/helper"
	"log"
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
			return c.JSON(helper.ErrorResponse("please upload the image"))
		}

		str := c.FormValue("camp_id")
		campID, err := strconv.Atoi(str)
		if err != nil {
			log.Println("Error in handler add", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		if err := ih.srv.Add(token, uint(campID), header); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(201, "success add new image"))
	}
}

func (ih *imageHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.Param("id")
		imageID, err := strconv.Atoi(str)
		if err != nil {
			log.Println("Error in handler delete", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		if err := ih.srv.Delete(token, uint(imageID)); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.NoContent(204)
	}
}
