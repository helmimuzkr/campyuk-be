package data

import (
	"campyuk-api/features/booking"

	"gorm.io/gorm"
)

type bookingQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) booking.BookingData {
	return &bookingQuery{
		db: db,
	}
}

func (bq *bookingQuery) Create(userID uint, virtualNumber string) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bq *bookingQuery) Update(bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bq *bookingQuery) List(userID uint) ([]booking.Core, error) {
	return []booking.Core{}, nil
}

func (bq *bookingQuery) GetByID(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}

func (bq *bookingQuery) Cancel(userID uint, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}
