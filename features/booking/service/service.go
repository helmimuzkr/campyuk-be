package service

import (
	"campyuk-api/features/booking"
	"campyuk-api/helper"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type bookingSrv struct {
	qry  booking.BookingData
	core coreapi.Client
}

func New(bd booking.BookingData, c coreapi.Client) booking.BookingService {
	return &bookingSrv{
		qry:  bd,
		core: c,
	}
}

func (bs *bookingSrv) Create(token interface{}, newBooking booking.Core) (booking.Core, error) {
	id, role := helper.ExtractToken(token)
	if role != "guest" {
		return booking.Core{}, errors.New("access is denied due to invalid credential")
	}

	// Assign some default transactions
	newBooking.Ticket = fmt.Sprintf("INV-%d-%s", id, time.Now().Format("20060102-150405"))
	newBooking.Status = "PENDING"
	newBooking.BookingDate = time.Now().Format("02-01-2006")

	// Charge midtrans
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newBooking.Ticket,
			GrossAmt: int64(newBooking.TotalPrice),
		},
		BankTransfer: &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(newBooking.Bank),
		},
		CustomExpiry: &coreapi.CustomExpiry{
			ExpiryDuration: 1,
			Unit:           "day",
		},
	}
	response, errMidtrans := bs.core.ChargeTransaction(chargeReq)
	if errMidtrans != nil {
		log.Println(errMidtrans)
		return booking.Core{}, errors.New("charge transaction failed due to internal server error")
	}
	newBooking.Bank = response.VaNumbers[0].Bank
	newBooking.VirtualNumber = response.VaNumbers[0].VANumber

	// Create booking
	res, err := bs.qry.Create(id, newBooking)
	if err != nil {
		log.Fatal(err.Error())
		return booking.Core{}, errors.New("internal server error")
	}

	return res, nil
}

func (bs *bookingSrv) Update(token interface{}, status string) error {
	return nil
}

func (bs *bookingSrv) List(token interface{}) ([]booking.Core, error) {
	return nil, nil
}

func (bs *bookingSrv) GetByID(token interface{}, bookingID uint) (booking.Core, error) {
	return booking.Core{}, nil
}
