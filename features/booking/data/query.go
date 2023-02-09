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

func (bd *bookingData) List(userID uint, role string, limit int, offset int) (int, []booking.Core, error) {
	var qryBooking, qryPagination, qryItem string

	if role == "host" {
		// Query for host
		qryBooking = "SELECT bookings.id, bookings.ticket, bookings.user_id, bookings.camp_id, camps.title,camps.address, camps.city, bookings.check_in, bookings.check_out, bookings.booking_date, bookings.total_price, bookings.status, 	bookings.bank, bookings.virtual_number FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE camps.host_id = ? ORDER BY bookings.id DESC LIMIT ? OFFSET ?"
		qryPagination = "SELECT COUNT(bookings.id) FROM bookings JOIN camps ON camps.id = bookings.camp_id WHERE camps.host_id = ?"
		qryItem = "SELECT items.name, items.price, rent_items.quantity, rent_items.cost FROM rent_items JOIN items ON items.id = rent_items.item_id JOIN camps ON camps.id = items.camp_id WHERE camps.host_id = ?"

	} else {
		// Query for guest
		qryBooking = "SELECT bookings.id, bookings.ticket, bookings.user_id, bookings.camp_id, camps.title,camps.address, camps.city, bookings.check_in, bookings.check_out, bookings.booking_date, bookings.total_price, bookings.status, bookings.bank, bookings.virtual_number FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE bookings.user_id = ? ORDER BY bookings.id DESC LIMIT ? OFFSET ?"
		qryPagination = "SELECT COUNT(id) FROM bookings WHERE user_id = ?"
		qryItem = "SELECT items.name, items.price, rent_items.quantity, rent_items.cost FROM rent_items JOIN items ON items.id = rent_items.item_id JOIN bookings ON bookings.id = rent_items.booking_id WHERE bookings.user_id = ?"
	}

	var models []BookingCamp
	tx := bd.db.Raw(qryBooking, userID, limit, offset).Find(&models)
	if tx.Error != nil {
		return 0, nil, tx.Error
	}

	for i := range models {
		tx = tx.Raw("SELECT image FROM images WHERE camp_id = ? ORDER BY id ASC", models[i].CampID).First(&models[i].Image)
		if tx.Error != nil {
			log.Println(tx.Error)
		}
	}

	for i := range models {
		tx = tx.Raw(qryItem, userID).Find(&models[i].Items)
		if tx.Error != nil {
			log.Println(tx.Error)
		}
	}

	var totalRecord int64
	tx = tx.Raw(qryPagination, userID).Find(&totalRecord)
	if tx.Error != nil {
		return 0, nil, tx.Error
	}

	return int(totalRecord), ToListCore(models), nil
}

func (bd *bookingData) GetByID(userID uint, bookingID uint, role string) (booking.Core, error) {
	var qryBooking, qryItem string

	if role == "host" {
		// Query for host
		qryBooking = "SELECT bookings.id, bookings.ticket, bookings.user_id, bookings.camp_id, camps.title, camps.latitude, camps.longitude, camps.address, camps.city, camps.price, bookings.check_in, bookings.check_out, bookings.booking_date, bookings.guest, bookings.camp_cost, bookings.total_price, bookings.status, bookings.bank, bookings.virtual_number FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE camps.host_id = ? AND bookings.id = ?"
		qryItem = "SELECT items.name, items.price, rent_items.quantity, rent_items.cost FROM rent_items JOIN items ON items.id = rent_items.item_id JOIN camps ON camps.id = items.camp_id WHERE camps.host_id = ? AND rent_items.booking_id = ?"

	} else {
		// Query for guest
		qryBooking = "SELECT bookings.id, bookings.ticket, bookings.user_id, bookings.camp_id, camps.title, camps.latitude, camps.longitude, camps.address, camps.city, camps.price, bookings.check_in, bookings.check_out, bookings.booking_date, bookings.guest, bookings.camp_cost, bookings.total_price, bookings.status, bookings.bank, bookings.virtual_number FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE bookings.user_id = ? AND bookings.id = ?"
		qryItem = "SELECT items.name, items.price, rent_items.quantity, rent_items.cost FROM rent_items JOIN items ON items.id = rent_items.item_id JOIN bookings ON bookings.id = rent_items.booking_id WHERE bookings.user_id = ? AND rent_items.booking_id = ? "
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

func (bd *bookingData) Update(userID uint, role string, bookingID uint, status string) error {
	var qry string
	if role == "host" {
		// role host
		qry = "UPDATE bookings JOIN camps ON camps.id = bookings.camp_id SET bookings.status = ? WHERE camps.host_id = ? AND bookings.id = ?"
		err := bd.decrementStock(bookingID)
		if err != nil {
			return err
		}

	} else if role == "guest" {
		// role guest
		qry = "UPDATE bookings SET status = ? WHERE user_id = ? AND id = ?"
	}

	tx := bd.db.Exec(qry, status, userID, bookingID)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (bd *bookingData) Callback(ticket string, status string) error {
	tx := bd.db.Model(&Booking{}).Where("ticket = ?", ticket).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}

	if status == "SUCCESS" {
		var bookingID uint
		tx = bd.db.Raw("SELECT id FROM bookings WHERE ticket = ?").First(&bookingID)
		if tx.Error != nil {
			return tx.Error
		}
		err := bd.decrementStock(bookingID)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (bd *bookingData) CreateEvent(bookingID uint) (booking.Core, error) {
	model := BookingCamp{}
	query := "SELECT bookings.id, users.email, bookings.ticket, camps.title, camps.latitude, camps.longitude, camps.address, camps.city, bookings.check_in, bookings.check_out, bookings.guest FROM bookings JOIN users ON users.id = bookings.user_id JOIN camps ON camps.id = bookings.camp_id WHERE bookings.id = ?"
	tx := bd.db.Raw(query, bookingID).First(&model)
	if tx.Error != nil {
		return booking.Core{}, tx.Error
	}

	return ToCore(model), nil
}

func (bd *bookingData) decrementStock(bookingID uint) error {
	itm := []RentItem{}
	tx := bd.db.Where("booking_id = ?", bookingID).Find(&itm)
	if tx.Error != nil {
		return tx.Error
	}

	for _, v := range itm {
		var stock int
		tx = bd.db.Raw("SELECT stock FROM items WHERE id = ?", v.ItemID).First(&stock)
		if stock < v.Quantity {
			return errors.New("stock not available")
		}
		tx = bd.db.Exec("UPDATE items SET stock = stock - ? WHERE id = ?", v.Quantity, v.ItemID)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}
