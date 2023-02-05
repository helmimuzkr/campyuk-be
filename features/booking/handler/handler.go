package handler

import (
	"campyuk-api/features/booking"
	"campyuk-api/helper"

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

func (bc *bookingController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (bc *bookingController) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (bc *bookingController) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (bc *bookingController) Accept() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (bc *bookingController) Cancel() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (bc *bookingController) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
