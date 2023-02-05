package handler

import (
	"campyuk-api/features/camp"
	"campyuk-api/helper"
	"strconv"

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

		return c.JSON(helper.SuccessResponse(201, "success add new camp"))
	}
}

func (ch *campHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		if token == nil {
			token = jwt.New(jwt.SigningMethodES256)
		}

		str := c.QueryParam("page")
		page, _ := strconv.Atoi(str)

		paginate, res, err := ch.srv.List(token, page)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		campResponse := []campResponse{}
		copier.Copy(&campResponse, &res)
		for i := range res {
			campResponse[i].Image = res[i].Images[0].ImageURL
		}

		pagination := paginationResponse{
			Page:        paginate["page"].(int),
			Limit:       paginate["limit"].(int),
			Offset:      paginate["offset"].(int),
			TotalRecord: paginate["totalRecord"].(int),
			TotalPage:   paginate["totalPage"].(int),
		}

		response := withPagination{
			Pagination: pagination,
			Data:       campResponse,
			Message:    "success show list camp",
		}

		return c.JSON(200, response)
	}
}

func (ch *campHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		if token == nil {
			token = jwt.New(jwt.SigningMethodES256)
		}

		str := c.Param("id")
		campID, _ := strconv.Atoi(str)

		res, err := ch.srv.GetByID(token, uint(campID))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := campDetailReponse{}
		copier.Copy(&response, &res)

		return c.JSON(helper.SuccessResponse(200, "success show detail camp", response))
	}
}

func (ch *campHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.Param("id")
		campID, _ := strconv.Atoi(str)

		cr := campRequest{}
		if err := c.Bind(&cr); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		document, _ := c.FormFile("document")

		updateCamp := camp.Core{}
		copier.Copy(&updateCamp, &cr)

		if err := ch.srv.Update(token, uint(campID), updateCamp, document); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(201, "success update camp"))
	}
}

func (ch *campHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		paramID := c.Param("id")
		campID, _ := strconv.Atoi(paramID)

		err := ch.srv.Delete(token, uint(campID))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success delete camp"))
	}
}

func (ch *campHandler) Accept() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		paramID := c.Param("id")
		campID, _ := strconv.Atoi(paramID)

		status := "ACCEPTED"

		err := ch.srv.RequestAdmin(token, uint(campID), status)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success accept camp"))
	}
}

func (ch *campHandler) Decline() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		paramID := c.Param("id")
		campID, _ := strconv.Atoi(paramID)

		status := "REJECTED"

		err := ch.srv.RequestAdmin(token, uint(campID), status)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success reject camp"))
	}
}
