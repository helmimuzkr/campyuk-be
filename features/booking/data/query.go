package data

import (
	"campyuk-api/features/booking"
	"errors"
	"log"

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

func (bd *bookingData) List(userID uint, role string, limit int, offset int) (int, []booking.Core, error) {
	model := []BookingCamp{}
	totalRecord := 0

	if role == "guest" {

	} else if role == "host" {

	} else {
		return 0, nil, errors.New("access is denied due to invalid credential")
	}

	return totalRecord, ToListCore(model), nil
}

func (bd *bookingData) GetByID(userID uint, bookingID uint, role string) (booking.Core, error) {
	var qryBooking, qryItem string

	if role == "guest" {
		qryBooking = "SELECT bookings.id, bookings.ticket, bookings.user_id, bookings.camp_id, camps.title, camps.latitude, camps.longitude, camps.address, camps.city, camps.price, bookings.check_in, bookings.check_out, bookings.booking_date, bookings.guest, bookings.camp_cost, bookings.total_price, bookings.status, bookings.bank, bookings.virtual_number FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE bookings.user_id = ? AND bookings.id = ?"
		qryItem = "SELECT items.name, items.price, rent_items.quantity, rent_items.cost FROM rent_items JOIN items ON items.id = rent_items.item_id JOIN bookings ON bookings.id = rent_items.booking_id WHERE bookings.user_id = ? AND rent_items.booking_id = ? "

	} else if role == "host" {
		qryBooking = "SELECT bookings.id, bookings.ticket, bookings.user_id, bookings.camp_id, camps.title, camps.latitude, camps.longitude, camps.address, camps.city, camps.price, bookings.check_in, bookings.check_out, bookings.booking_date, bookings.guest, bookings.camp_cost, bookings.total_price, bookings.status, bookings.bank, bookings.virtual_number FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE camps.host_id = ? AND bookings.id = ?"
		qryItem = "SELECT items.name, items.price, rent_items.quantity, rent_items.cost FROM rent_items JOIN items ON items.id = rent_items.item_id JOIN camps ON camps.id = items.camp_id WHERE camps.host_id = ? AND rent_items.booking_id = ?"

	} else {
		return booking.Core{}, errors.New("access is denied due to invalid credential")
	}

	var model BookingCamp
	tx := bd.db.Raw(qryBooking, userID, bookingID).Find(&model)
	if tx.Error != nil {
		return booking.Core{}, tx.Error
	}

	if model.ID <= 0 {
		return booking.Core{}, errors.New("booking order not found")
	}

	var image string
	tx = tx.Raw("SELECT image FROM images WHERE camp_id = ? ORDER BY id ASC", model.CampID).First(&image)
	if tx.Error != nil {
		log.Println(tx.Error)
	}

	itemModel := []Item{}
	tx = tx.Raw(qryItem, userID, bookingID).Find(&itemModel)
	if tx.Error != nil {
		log.Println(tx.Error)
	}

	model.Image = image
	model.Items = itemModel

	return ToCore(model), nil
}

func (bd *bookingData) Update(userID uint, bookingID uint, status string) error {
	// role host
	qry := "UPDATE bookings JOIN camps ON camps.id = bookings.camp_id SET bookings.status = ? WHERE camps.host_id = ? AND bookings.id = ?"
	tx := bd.db.Exec(qry, status, userID, bookingID)
	if tx.Error != nil {
		return tx.Error
	}

	// role guest
	qry2 := "UPDATE bookings SET status = ? WHERE user_id = ? AND id = ?"
	tx = bd.db.Exec(qry2, status, userID, bookingID)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (bd *bookingData) Callback(ticket string, status string) error {
	err := bd.db.Model(&Booking{}).Where("ticket = ?", ticket).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}
