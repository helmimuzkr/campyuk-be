package handler

import (
	"campyuk-api/features/booking"

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
		return nil
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
