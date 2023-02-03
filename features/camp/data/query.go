package data

import (
	"campyuk-api/features/camp"

	"gorm.io/gorm"
)

type campData struct {
	db *gorm.DB
}

func New(db *gorm.DB) camp.CampData {
	return &campData{db: db}
}

func (cd *campData) Add(userID uint, newCamp camp.Core) error {
	// Create camp
	cm := ToData(userID, newCamp)
	tx := cd.db.Create(&cm)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	// Inserting image to camp
	cim := ToImageData(cm.ID, newCamp.Images)
	for _, v := range cim {
		// Kenapa pakai exec ketimbang batch create dari gorm? dikarenakan kena panic error.
		tx = tx.Exec("INSERT INTO camp_images(camp_id, image) VALUES(?, ?)", v.CampID, v.Image)
		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	tx.Commit()

	return nil
}

func (cd *campData) List(userID uint, role string) ([]camp.Core, error) {
	var cm []CampModel

	switch role {
	case "host":
		res, err := cd.listCampHost(userID)
		if err != nil {
			return nil, err
		}
		cm = res
	case "admin":
		res, err := cd.listCampAdmin()
		if err != nil {
			return nil, err
		}
		cm = res
	default:
		res, err := cd.listCampUser()
		if err != nil {
			return nil, err
		}
		cm = res
	}

	return ToListCampCore(cm), nil
}

func (cd *campData) GetByID(userID uint, role string, campID uint) (camp.Core, error) {
	return camp.Core{}, nil
}
func (cd *campData) Update(userID uint, campID uint, updateCamp camp.Core) error {
	return nil
}
func (cd *campData) Delete(userID uint, campID uint) error {
	return nil
}
func (cd *campData) RequestAdmin(userID uint, campID uint) error {
	return nil
}

// ------------------------
// Functions that are not include in the contract
// ------------------------

func (cd *campData) listCampUser() ([]CampModel, error) {
	cm := []CampModel{}
	// Select camp
	qc := "SELECT camps.id, camps.verification_status, users.fullname, camps.title, camps.price, camps.distance, camps.city FROM camps JOIN users ON users.id = camps.host_id WHERE camps.verification_status = 'ACCEPTED'"
	tx := cd.db.Raw(qc).Find(&cm)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Find camp image
	for i := range cm {
		ci := []CampImage{}
		tx = tx.Raw("SELECT image FROM camp_images WHERE camp_id = ?", cm[i].ID).Find(&ci)
		if tx.Error != nil {
			return nil, tx.Error
		}
		cm[i].CampImages = ci
	}

	return cm, nil
}

func (cd *campData) listCampHost(userID uint) ([]CampModel, error) {
	cm := []CampModel{}
	// Select camp
	qc := "SELECT camps.id, camps.verification_status, users.fullname, camps.title, camps.price, camps.distance,camps.city FROM camps JOIN users ON users.id = camps.host_id WHERE users.id = ?"
	tx := cd.db.Raw(qc, userID).Find(&cm)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Find camp image
	for i := range cm {
		ci := []CampImage{}
		tx = tx.Raw("SELECT image FROM camp_images WHERE camp_id = ?", cm[i].ID).Find(&ci)
		if tx.Error != nil {
			return nil, tx.Error
		}
		cm[i].CampImages = ci
	}

	return cm, nil
}

func (cd *campData) listCampAdmin() ([]CampModel, error) {
	cm := []CampModel{}
	// Select camp
	qc := "SELECT camps.id, camps.verification_status, users.fullname, camps.title, camps.price, camps.distance,camps.city FROM camps JOIN users ON users.id = camps.host_id WHERE camps.verification_status = 'PENDING'"
	tx := cd.db.Raw(qc).Find(&cm)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Find camp image
	for i := range cm {
		ci := []CampImage{}
		tx = tx.Raw("SELECT image FROM camp_images WHERE camp_id = ?", cm[i].ID).Find(&ci)
		if tx.Error != nil {
			return nil, tx.Error
		}
		cm[i].CampImages = ci
	}

	return cm, nil
}

// qc := "SELECT camps.id, camps.verification_status, users.fullname, camps.title, camps.price, camps.description, camps.latitude, camps.longitude, camps.distance, camps.address, camps.city, camps.document FROM camps JOIN users ON users.id = camps.host_id WHERE camps.host_id = ?"
