package service

import (
	"campyuk-api/features/booking"
	"campyuk-api/helper"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type bookingSrv struct {
	qry     booking.BookingData
	payment helper.PaymentGateway
}

func New(bd booking.BookingData, p helper.PaymentGateway) booking.BookingService {
	return &bookingSrv{
		qry:     bd,
		payment: p,
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

	// Charge transaction to midtrans and get the response
	vaNumber, errMidtrans := bs.payment.ChargeTransaction(newBooking.Ticket, newBooking.TotalPrice, newBooking.Bank)
	if errMidtrans != nil {
		log.Println(errMidtrans)
		return booking.Core{}, errors.New("charge transaction failed due to internal server error, please try again")
	}
	newBooking.VirtualNumber = vaNumber

	// Create booking
	res, err := bs.qry.Create(id, newBooking)
	if err != nil {
		log.Println(err.Error())
		return booking.Core{}, errors.New("internal server error")
	}

	return res, nil
}

func (bs *bookingSrv) List(token interface{}, page int) (map[string]interface{}, []booking.Core, error) {
	userID, role := helper.ExtractToken(token)
	if role != "guest" && role != "host" {
		return nil, nil, errors.New("access is denied due to invalid credential")
	}

	if page < 1 {
		page = 1
	}
	limit := 4
	// Calculate offset
	offset := (page - 1) * limit

	totalRecord, res, err := bs.qry.List(userID, role, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal server error")
	}

	totalPage := int(math.Ceil(float64(totalRecord) / float64(limit)))

	pagination := make(map[string]interface{})
	pagination["page"] = page
	pagination["limit"] = limit
	pagination["offset"] = offset
	pagination["totalRecord"] = totalRecord
	pagination["totalPage"] = totalPage

	return pagination, res, nil
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
		} else if strings.Contains(err.Error(), "stock") {
			msg = err.Error()
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

	status = strings.ToUpper(status)

	err := bs.qry.Callback(ticket, status)
	if err != nil {
		log.Println("callback error", err)
		return errors.New("internal server error")
	}

	return nil
}
