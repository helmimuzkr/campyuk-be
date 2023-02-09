package service

import (
	"campyuk-api/features/booking"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func setupTest(t *testing.T) (*mocks.PaymentGateway, *mocks.BookingData, booking.BookingService, *mocks.GoogleAPI) {
	midtrans := mocks.NewPaymentGateway(t)
	data := mocks.NewBookingData(t)
	google := mocks.NewGoogleAPI(t)
	srv := New(data, midtrans, google)

	return midtrans, data, srv, google
}

func dataSample() (booking.Core, booking.Core) {
	inData := booking.Core{
		Ticket: fmt.Sprintf("INV-%d-%s", uint(1), time.Now().Format("20060102-150405")),
		CampID: uint(1),
		Items: []booking.Item{
			{ID: uint(2), Quantity: 1, RentCost: 50000},
		},
		CheckIn:       "2023-02-09",
		CheckOut:      "2023-02-10",
		Guest:         1,
		CampCost:      20000,
		TotalPrice:    70000,
		Status:        "PENDING",
		BookingDate:   time.Now().Format("2006-01-02"),
		Bank:          "bca",
		VirtualNumber: "90316950939",
	}

	resData := booking.Core{
		ID:        uint(1),
		Ticket:    fmt.Sprintf("INV-%d-%s", uint(1), time.Now().Format("20060102-150405")),
		CampID:    uint(1),
		Title:     "Tanakita camp",
		CampPrice: 20000,
		Image:     "https://res.cloudinary.com/dnji8pgyl/image/upload/v1675500679/file/20230204-155113.webp",
		Latitude:  -6.208987101998694,
		Longitude: 106.79970296358913,
		Address:   "Jl. Spjljk",
		City:      "Pochinki",
		Items: []booking.Item{
			{Name: "Small tent", Price: 500000, Quantity: 1, RentCost: 50000},
		},
		CheckIn:       "2023-02-09",
		CheckOut:      "2023-02-10",
		Guest:         1,
		CampCost:      20000,
		TotalPrice:    70000,
		Status:        "PENDING",
		BookingDate:   time.Now().Format("2006-01-02"),
		Bank:          "bca",
		VirtualNumber: "90316950939",
	}

	return inData, resData
}

func TestCreateBooking(t *testing.T) {
	midtrans, data, srv, _ := setupTest(t)
	inData, resData := dataSample()

	t.Run("Succress create new order", func(t *testing.T) {
		midtrans.On("ChargeTransaction", inData.Ticket, inData.TotalPrice, inData.Bank).Return("90316950939", nil).Once()

		data.On("Create", uint(1), inData).Return(booking.Core{ID: uint(1)}, nil).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, inData)

		assert.Nil(t, err)
		assert.Equal(t, resData.ID, actual.ID)
		data.AssertExpectations(t)
	})

	t.Run("Error access is denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, inData)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
		assert.Empty(t, actual)
		data.AssertExpectations(t)
	})

	t.Run("Charge failed", func(t *testing.T) {
		midtrans.On("ChargeTransaction", inData.Ticket, inData.TotalPrice, inData.Bank).Return("", errors.New("charge transaction failed due to internal server error")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, inData)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "charge transaction failed due to internal server error")
		assert.Empty(t, actual)
		data.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		midtrans.On("ChargeTransaction", inData.Ticket, inData.TotalPrice, inData.Bank).Return("90316950939", nil).Once()
		data.On("Create", uint(1), inData).Return(booking.Core{}, errors.New("internal server error")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, inData)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Empty(t, actual)
		data.AssertExpectations(t)
	})
}

func TestList(t *testing.T) {
	_, data, srv, _ := setupTest(t)
	_, resData := dataSample()
	listData := []booking.Core{resData, resData}
	listData[0].ID = 1
	listData[1].ID = 2

	t.Run("Succress display order list", func(t *testing.T) {
		data.On("List", uint(1), "guest", 4, 0).Return(2, listData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		pagination, actual, err := srv.List(token, 1)

		assert.Nil(t, err)
		assert.Equal(t, listData[1].ID, actual[1].ID)
		assert.Equal(t, 2, pagination["totalRecord"])
		data.AssertExpectations(t)
	})

	t.Run("Error access is denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "")
		token := tkn.(*jwt.Token)
		token.Valid = false

		pagination, actual, err := srv.List(token, 1)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
		assert.Nil(t, actual)
		assert.Nil(t, pagination)
	})

	t.Run("Database error", func(t *testing.T) {
		data.On("List", uint(1), "guest", 4, 0).Return(0, nil, errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		pagination, actual, err := srv.List(token, 1)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Nil(t, actual)
		assert.Nil(t, pagination)
		data.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	_, data, srv, _ := setupTest(t)
	_, resData := dataSample()

	t.Run("Succress show booking detail", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1), "guest").Return(resData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetByID(token, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, resData.ID, actual.ID)
		assert.NotNil(t, actual.Items)
		data.AssertExpectations(t)
	})

	t.Run("Error access is denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "")
		token := tkn.(*jwt.Token)
		token.Valid = false

		actual, err := srv.GetByID(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
		assert.Empty(t, actual)
		assert.Nil(t, actual.Items)
	})

	t.Run("Database error", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1), "guest").Return(booking.Core{}, errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetByID(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Empty(t, actual)
		assert.Nil(t, actual.Items)
		data.AssertExpectations(t)
	})

	t.Run("Booking not found", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1), "guest").Return(booking.Core{}, errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetByID(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "booking order not found")
		assert.Empty(t, actual)
		assert.Nil(t, actual.Items)
		data.AssertExpectations(t)
	})
}

func TestAccept(t *testing.T) {
	_, data, srv, _ := setupTest(t)

	t.Run("Succress accept order", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "SUCCESS").Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Accept(token, uint(1), "SUCCESS")

		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("Error access is denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "")
		token := tkn.(*jwt.Token)
		token.Valid = false

		err := srv.Accept(token, uint(1), "SUCCESS")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
	})

	t.Run("Database error", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "SUCCESS").Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Accept(token, uint(1), "SUCCESS")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})

	t.Run("Booking not found", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "SUCCESS").Return(errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Accept(token, uint(1), "SUCCESS")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "booking not found")
		data.AssertExpectations(t)
	})

	t.Run("Stock not available", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "SUCCESS").Return(errors.New("stock not available")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Accept(token, uint(1), "SUCCESS")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "stock not available")
		data.AssertExpectations(t)
	})
}

func TestCancel(t *testing.T) {
	_, data, srv, _ := setupTest(t)

	t.Run("Succress cancel order", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "CANCEL").Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1), "CANCEL")

		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("Error access is denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "")
		token := tkn.(*jwt.Token)
		token.Valid = false

		err := srv.Cancel(token, uint(1), "CANCEL")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
	})

	t.Run("Database error", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "CANCEL").Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1), "CANCEL")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})

	t.Run("Booking not found", func(t *testing.T) {
		data.On("Update", uint(1), "host", uint(1), "CANCEL").Return(errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1), "CANCEL")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "booking not found")
		data.AssertExpectations(t)
	})
}

func TestCallback(t *testing.T) {
	_, data, srv, _ := setupTest(t)
	inData, _ := dataSample()

	t.Run("Succress callback ", func(t *testing.T) {
		data.On("Callback", inData.Ticket, "SUCCESS").Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Callback(inData.Ticket, "settlement")

		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		data.On("Callback", inData.Ticket, "SUCCESS").Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Callback(inData.Ticket, "settlement")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestCreateEvent(t *testing.T) {
	_, data, srv, google := setupTest(t)
	_, resData := dataSample()

	startTime, _ := time.Parse("2006-01-02", resData.CheckIn)

	endTime, _ := time.Parse("2006-01-02", resData.CheckOut)

	startRFC := startTime.Format(time.RFC3339)
	endRFC := endTime.Format(time.RFC3339)

	resGoogle := &oauth2.Token{}
	detailCal := helper.CalendarDetail{
		Summay:   "Camping",
		Location: resData.Address,
		Start:    startRFC,
		End:      endRFC,
		// nanti diisi email guest dan host
		Emails: []string{resData.Email}, // email guest
	}

	t.Run("Success create event in calendar", func(t *testing.T) {
		google.On("GetToken", "code").Return(resGoogle, nil).Once()
		data.On("CreateEvent", resData.ID).Return(resData, nil).Once()
		google.On("CreateCalendar", resGoogle, detailCal).Return(nil).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.Nil(t, err)
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})
	t.Run("Failed to get token", func(t *testing.T) {
		google.On("GetToken", "code").Return(nil, errors.New("failed")).Once()
		// data.On("CreateEvent", resData.ID).Return(resData, nil).Once()
		// google.On("CreateCalendar", resGoogle, detailCal).Return(nil).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to create event in calendar")
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})

	t.Run("booking not found", func(t *testing.T) {
		google.On("GetToken", "code").Return(resGoogle, nil).Once()
		data.On("CreateEvent", resData.ID).Return(booking.Core{}, errors.New("not found")).Once()
		// google.On("CreateCalendar", resGoogle, detailCal).Return(nil).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "booking not found")
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		google.On("GetToken", "code").Return(resGoogle, nil).Once()
		data.On("CreateEvent", resData.ID).Return(booking.Core{}, errors.New("query error")).Once()
		// google.On("CreateCalendar", resGoogle, detailCal).Return(nil).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})

	t.Run("parsing start time error", func(t *testing.T) {
		google.On("GetToken", "code").Return(resGoogle, nil).Once()

		resData.CheckIn = ""
		data.On("CreateEvent", resData.ID).Return(resData, nil).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to create event in calendar")
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})

	t.Run("parsing ebd time error", func(t *testing.T) {
		google.On("GetToken", "code").Return(resGoogle, nil).Once()

		resData.CheckOut = ""
		data.On("CreateEvent", resData.ID).Return(resData, nil).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to create event in calendar")
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})

	t.Run("create calendar failed", func(t *testing.T) {
		google.On("GetToken", "code").Return(resGoogle, nil).Once()
		data.On("CreateEvent", resData.ID).Return(resData, nil).Once()
		google.On("CreateCalendar", resGoogle, detailCal).Return(errors.New("failed to create")).Once()

		err := srv.CreateEvent("code", resData.ID)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to create event in calendar")
		data.AssertExpectations(t)
		google.AssertExpectations(t)
	})
}
