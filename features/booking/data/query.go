package data

import (
	"campyuk-api/features/booking"

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

func (bq *bookingData) Create(userID uint, virtualNumber string) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bq *bookingData) Update(bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bq *bookingData) List(userID uint) ([]booking.Core, error) {
	return []booking.Core{}, nil
}

func (bq *bookingData) GetByID(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bq *bookingData) Cancel(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}
