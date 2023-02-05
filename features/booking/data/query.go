package data

import (
	"campyuk-api/features/booking"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

type bookingData struct {
	db *gorm.DB
}

func New(db *gorm.DB) booking.BookingData {
	return &bookingData{
		db: db,
	}
}

func (bd *bookingData) Create(userID uint, newBooking booking.Core) (booking.Core, error) {
	model := ToData(userID, newBooking)
	tx := bd.db.Create(&model)
	if tx.Error != nil {
		return booking.Core{}, tx.Error
	}

	return booking.Core{ID: model.ID}, nil
}

func (bd *bookingData) Update(userID uint, role string, bookingID uint, status string) error {
	return nil
}

func (bd *bookingData) List(userID uint) ([]booking.Core, error) {
	return []booking.Core{}, nil
}

func (bd *bookingData) GetByID(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

// ------------------------
// Functions that are not include in the contract
// ------------------------

func (bd *bookingData) ChargePayment(newBooking booking.Core, items *[]midtrans.ItemDetails) (*coreapi.ChargeResponse, error) {

	return nil, nil
}
