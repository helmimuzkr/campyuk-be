package booking

import "github.com/labstack/echo/v4"

type Core struct {
	ID            uint
	Ticket        string
	GuestID       uint   // Guest
	Email         string // Email guest
	CampID        uint
	Title         string
	Image         string
	Latitude      complex64
	Longitude     complex64
	Address       string
	City          string
	CampPrice     int
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
	List() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Accept() echo.HandlerFunc
	Cancel() echo.HandlerFunc
	Callback() echo.HandlerFunc
}

type BookingService interface {
	Create(token interface{}, newBooking Core) (Core, error)
	List(token interface{}, page int) (map[string]interface{}, []Core, error)
	GetByID(token interface{}, bookingID uint) (Core, error)
	Accept(token interface{}, bookingID uint, status string) error
	Cancel(token interface{}, bookingID uint, status string) error
	Callback(ticket string, status string) error
}

type BookingData interface {
	Create(userID uint, newBooking Core) (Core, error)
	Update(userID uint, role string, bookingID uint, status string) error
	List(userID uint, role string, limit int, offset int) (int, []Core, error)
	GetByID(userID uint, bookingID uint, role string) (Core, error)
	Callback(ticket string, status string) error
}
