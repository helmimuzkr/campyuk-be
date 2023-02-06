package booking

import "github.com/labstack/echo/v4"

type Core struct {
	ID            uint
	Ticket        string
	UserID        uint   // Guest
	Email         string // Email guest
	CampID        uint
	Title         string
	Image         string
	Address       string
	City          string
	CampPrice     string
	CheckIn       string
	CheckOut      string
	BookingDate   string
	Guest         int
	CampCost      int
	Items         []Item
	TotalPrice    int
	Status        string
	Bank          string
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
	Update() echo.HandlerFunc
	List() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Accept() echo.HandlerFunc
	Cancel() echo.HandlerFunc
	Callback() echo.HandlerFunc
}

type BookingService interface {
	Create(token interface{}, newBooking Core) (Core, error)
	Update(token interface{}, updateBooking Core) error
	List(token interface{}) ([]Core, error)
	GetByID(token interface{}, bookingID uint) (Core, error)
	Callback(ticket string, status string) error
	RequestHost(token interface{}, bookingID uint, status string) error
}

type BookingData interface {
	Create(userID uint, newBooking Core) (Core, error)
	Update(userID uint, role string, updateBooking Core) error
	List(userID uint) ([]Core, error)
	GetByID(userID uint, bookingID uint) (Core, error)
	Callback(ticket string, status string) error
	RequestHost(bookingID uint, status string) error
}
