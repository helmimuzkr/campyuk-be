package data

import (
	"campyuk-api/features/booking"
	"strings"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID        uint // Guest
	CampID        uint
	Ticket        string
	CheckIn       string
	CheckOut      string
	BookingDate   string
	Guest         int
	CampCost      int
	TotalPrice    int
	Status        string
	Bank          string
	VirtualNumber string
	RentItem      []RentItem `gorm:"foreignKey:BookingID"`
}

type RentItem struct {
	gorm.Model
	BookingID uint
	ItemID    uint
	Quantity  int
	Cost      int
}

type BookingCamp struct {
	ID            uint
	Ticket        string
	UserID        uint   // Guest
	Email         string // Guest
	CampID        uint
	Title         string
	Image         string
	Latitude      complex64
	Longitude     complex64
	Address       string
	City          string
	CampPrice     int `gorm:"column:price"`
	CheckIn       string
	CheckOut      string
	BookingDate   string
	Guest         int
	CampCost      int
	TotalPrice    int
	Status        string
	Bank          string
	VirtualNumber string
	Items         []Item `gorm:"-"`
}

type Item struct {
	Name     string
	Price    int
	Quantity int
	RentCost int `gorm:"column:cost"`
}

func ToData(userID uint, core booking.Core) Booking {
	ri := []RentItem{}
	for _, v := range core.Items {
		ri = append(ri, RentItem{
			ItemID:    v.ID,
			BookingID: core.ID,
			Quantity:  v.Quantity,
			Cost:      v.RentCost,
		})
	}

	b := Booking{
		Model:         gorm.Model{ID: core.ID},
		UserID:        userID,
		CampID:        core.CampID,
		Ticket:        core.Ticket,
		CheckIn:       core.CheckIn,
		CheckOut:      core.CheckOut,
		BookingDate:   core.BookingDate,
		Guest:         core.Guest,
		CampCost:      core.CampCost,
		TotalPrice:    core.TotalPrice,
		Status:        core.Status,
		Bank:          core.Bank,
		VirtualNumber: core.VirtualNumber,
		RentItem:      ri,
	}

	return b
}

func ToCore(data BookingCamp) booking.Core {
	return booking.Core{
		ID:            data.ID,
		GuestID:       data.UserID,
		Email:         data.Email,
		Ticket:        data.Ticket,
		Title:         data.Title,
		Image:         data.Image,
		Latitude:      data.Latitude,
		Longitude:     data.Longitude,
		Address:       data.Address,
		City:          data.City,
		CampPrice:     data.CampPrice,
		CheckIn:       data.CheckIn,
		CheckOut:      data.CheckOut,
		BookingDate:   data.BookingDate,
		Guest:         data.Guest,
		CampCost:      data.CampCost,
		TotalPrice:    data.TotalPrice,
		Status:        strings.ToUpper(data.Status),
		Bank:          data.Bank,
		VirtualNumber: data.VirtualNumber,
		Items:         ToItemsCore(data.Items),
	}
}

func ToItemsCore(data []Item) []booking.Item {
	var cores []booking.Item

	for _, v := range data {
		c := booking.Item{}
		c.Name = v.Name
		c.Price = v.Price
		c.Quantity = v.Quantity
		c.Price = v.Price
		c.RentCost = v.RentCost

		cores = append(cores, c)
	}

	return cores
}

func ToListCore(data []BookingCamp) []booking.Core {
	cores := []booking.Core{}

	for _, v := range data {
		cores = append(cores, ToCore(v))
	}

	return cores
}
