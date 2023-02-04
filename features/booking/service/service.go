package service

import "campyuk-api/features/booking"

type bookingSrv struct {
	qry booking.BookingData
}

func New(bd booking.BookingData) booking.BookingService {
	return &bookingSrv{
		qry: bd,
	}
}

func (bs *bookingSrv) Create(token interface{}) (booking.Core, error)

func (bs *bookingSrv) Update(token interface{}) (booking.Core, error)

func (bs *bookingSrv) List(token interface{}) ([]booking.Core, error)

func (bs *bookingSrv) GetByID(token interface{}, bookingID uint) (booking.Core, error)

func (bs *bookingSrv) Cancel(token interface{}, bookingID uint) (booking.Core, error)
