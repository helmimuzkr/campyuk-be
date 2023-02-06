package service

import (
	"campyuk-api/features/booking"
	"campyuk-api/helper"
	"errors"
	"fmt"
	"log"
	"strings"
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

func (bs *bookingSrv) List(token interface{}, page int) ([]booking.Core, error) {

	return nil, nil
}

func (bs *bookingSrv) GetByID(token interface{}, bookingID uint) (booking.Core, error) {
	userID, role := helper.ExtractToken(token)
	if role != "guest" && role != "host" {
		return booking.Core{}, errors.New("access is denied due to invalid credential")
	}

	res, err := bs.qry.GetByID(userID, bookingID, role)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "access is denied") {
			return booking.Core{}, err
		}
		if strings.Contains(err.Error(), "not found") {
			return booking.Core{}, errors.New("booking order not found")
		}

		return booking.Core{}, errors.New("internal server error")
	}

	return res, nil
}

func (bs *bookingSrv) Accept(token interface{}, bookingID uint, status string) error {
	id, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	if err := bs.qry.Update(id, role, bookingID, status); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "booking not found"
		} else {
			msg = "internal server errorr"
		}
		return errors.New(msg)
	}

	return nil
}

func (bs *bookingSrv) Cancel(token interface{}, bookingID uint, status string) error {
	id, role := helper.ExtractToken(token)
	if role != "guest" && role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	if err := bs.qry.Update(id, role, bookingID, status); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "booking not found"
		} else {
			msg = "internal server errorr"
		}
		return errors.New(msg)
	}

	return nil
}

func (bs *bookingSrv) Callback(ticket string, status string) error {
	if status == "settlement" {
		status = "SUCCESS"
	}
	err := bs.qry.Callback(ticket, status)
	if err != nil {
		log.Println("callback error", err)
		return errors.New("internal server error")
	}

	return nil
}
