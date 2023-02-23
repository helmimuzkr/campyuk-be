package booking

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID            uint
	Ticket        string
	GuestID       uint   // Guest
	Email         string // Email guest
	CampID        uint
	Title         string
	Image         string
	Latitude      float64
	Longitude     float64
	Address       string
	City          string
	CampPrice     int
	CheckIn       string `validate:"required"`
	CheckOut      string `validate:"required"`
	BookingDate   string
	Guest         int `validate:"required"`
	CampCost      int
	Items         []Item
	TotalPrice    int `validate:"required"`
	Status        string
	Bank          string `validate:"required"`
	VirtualNumber string
}

type Item struct {
	ID       uint
	Name     string
	Price    int
	Quantity int
	RentCost int
}

type BookingHandler interface {
	Create() echo.HandlerFunc
	List() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Accept() echo.HandlerFunc
	Cancel() echo.HandlerFunc
	Callback() echo.HandlerFunc
	CreateReminder() echo.HandlerFunc
}

type BookingService interface {
	Create(token interface{}, newBooking Core) (Core, error)
	List(token interface{}, page int) (map[string]interface{}, []Core, error)
	GetByID(token interface{}, bookingID uint) (Core, error)
	Accept(token interface{}, bookingID uint, status string) error
	Cancel(token interface{}, bookingID uint, status string) error
	CreateReminder(token interface{}, bookingID uint) (string, error)
	Callback(ticket string, status string) error
}

type BookingRepository interface {
	Create(userID uint, newBooking Core) (Core, error)
	Update(userID uint, role string, bookingID uint, status string) error
	List(userID uint, role string, limit int, offset int) (int, []Core, error)
	GetByID(userID uint, bookingID uint, role string) (Core, error)
	Callback(ticket string, status string) error
}

type GoogleGateway interface {
	CreateEvent(detailEvent map[string]string) (string, error)
}

type PaymentGateway interface {
	ChargeTransaction(orderID string, grossAmt int, bank string) (string, error)
}
