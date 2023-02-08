package handler

import (
	"campyuk-api/features/booking"
	"campyuk-api/helper"
	"log"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type bookingController struct {
	srv booking.BookingService
}

func New(bs booking.BookingService) booking.BookingHandler {
	return &bookingController{
		srv: bs,
	}
}

func (bc *bookingController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		br := bookingRequest{}
		if err := c.Bind(&br); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		newBooking := booking.Core{}
		copier.Copy(&newBooking, &br)
		if err := copier.Copy(&newBooking, &br); err != nil {
			log.Println("create booking", err)
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		res, err := bc.srv.Create(token, newBooking)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := struct {
			BookingID uint `json:"booking_id"`
		}{
			BookingID: res.ID,
		}

		return c.JSON(helper.SuccessResponse(201, "success booking", response))
	}
}

func (bc *bookingController) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.QueryParam("page")
		page, err := strconv.Atoi(str)
		if err != nil {
			log.Println("convert id error", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		paginate, res, err := bc.srv.List(token, page)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		bookingResp := []bookingResponse{}
		copier.Copy(&bookingResp, &res)

		pagination := helper.PaginationResponse{
			Page:        paginate["page"].(int),
			Limit:       paginate["limit"].(int),
			Offset:      paginate["offset"].(int),
			TotalRecord: paginate["totalRecord"].(int),
			TotalPage:   paginate["totalPage"].(int),
		}

		webResponse := helper.WithPagination{
			Pagination: pagination,
			Data:       bookingResp,
			Message:    "show all transaction success",
		}

		return c.JSON(200, webResponse)
	}
}

func (bc *bookingController) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.Param("id")
		bookingID, err := strconv.Atoi(str)
		if err != nil {
			log.Println("convert id error", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		res, err := bc.srv.GetByID(token, uint(bookingID))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := bookingDetailResponse{}
		copier.Copy(&response, &res)

		return c.JSON(helper.SuccessResponse(200, "success show detail booking", response))
	}
}

func (bc *bookingController) Accept() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		paramID := c.Param("id")
		bookingID, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("convert id error", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		status := "SUCCESS"

		err = bc.srv.Accept(token, uint(bookingID), status)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success accept booking"))
	}
}

func (bc *bookingController) Cancel() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		paramID := c.Param("id")
		bookingID, err := strconv.Atoi(paramID)
		if err != nil {
			log.Println("convert id error", err.Error())
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		status := "CANCELLED"

		err = bc.srv.Cancel(token, uint(bookingID), status)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success cancel booking"))
	}
}

func (bc *bookingController) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		cb := Callback{}
		if err := c.Bind(&cb); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		err := bc.srv.Callback(cb.OrderID, cb.TransactionStatus)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success update transaction"))
	}
}
