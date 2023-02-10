package handler

import (
	"campyuk-api/features/item"
	"campyuk-api/helper"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type itemHandler struct {
	srv item.ItemService
}

func New(is item.ItemService) item.ItemHandler {
	return &itemHandler{
		srv: is,
	}
}

func (ih *itemHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := AddItemRequest{}
		err := c.Bind(&input)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "wrong input format"})
		}

		res, err := ih.srv.Add(c.Get("user"), input.CampID, *RequestToCore(input))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		log.Println(res)
		return c.JSON(helper.SuccessResponse(http.StatusCreated, "success add new item"))
	}
}

func (ih *itemHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		paramID := c.Param("id")
		itemID, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("convert id error", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		input := UpdateItemRequest{}
		err = c.Bind(&input)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "wrong input format"})
		}

		res, err := ih.srv.Update(c.Get("user"), uint(itemID), *RequestToCore(input))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		log.Println(res)
		return c.JSON(helper.SuccessResponse(http.StatusOK, "success update item"))
	}
}

func (ih *itemHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		paramID := c.Param("id")
		itemID, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("convert id error", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		err = ih.srv.Delete(c.Get("user"), uint(itemID))
		if err != nil {
			log.Println("trouble :  ", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.NoContent(204)
	}
}
