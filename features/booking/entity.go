package booking

import "github.com/labstack/echo/v4"

type Core struct {
	ID            uint
	Checkin       string
	Checkout      string
	BookingDate   string
	Guest         int
	CampCost      int
	TotalPrice    int
	Status        string
	Bank          string
	VirtualNumber string
}

type BookingHandler interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	List() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Accept() echo.HandlerFunc
	Cancel() echo.HandlerFunc
}

type BookingService interface {
	Create(token interface{}) (Core, error)
	Update(token interface{}) (Core, error)
	List(token interface{}) ([]Core, error)
	GetByID(token interface{}, bookingID uint) (Core, error)
	Cancel(token interface{}, bookingID uint) (Core, error)
}

type BookingData interface {
	Create(userID uint, virtualNumber string) (Core, error)
	Update(bookingID uint) (Core, error)
	List(userID uint) ([]Core, error)
	GetByID(userID uint, bookingID uint) (Core, error)
	Cancel(userID uint, bookingID uint) (Core, error)
}
