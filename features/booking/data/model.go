package data

import (
	"campyuk-api/features/booking"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	Checkin       string
	Checkout      string
	BookingDate   string
	Guest         int
	CampCost      int
	TotalPrice    int
	Status        string
	Bank          string
	VirtualNumber string
	UserID        uint
	CampID        uint
}

func ToCore(data Booking) booking.Core {
	return booking.Core{
		ID:          data.ID,
		Checkin:     data.Checkin,
		Checkout:    data.Checkout,
		BookingDate: data.BookingDate,
		Guest:       data.Guest,
		CampCost:    data.CampCost,
		TotalPrice:  data.TotalPrice,
		Status:      data.Status,
	}
}

func CoreToData(data booking.Core) Booking {
	return Booking{
		Model:       gorm.Model{ID: data.ID},
		Checkin:     data.Checkin,
		Checkout:    data.Checkout,
		BookingDate: data.BookingDate,
		Guest:       data.Guest,
		CampCost:    data.CampCost,
		TotalPrice:  data.TotalPrice,
		Status:      data.Status,
	}
}
