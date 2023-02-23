package repository

import (
	"errors"
	"log"

	"campyuk-api/features/user"

	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserRepository {
	return &userQuery{
		db: db,
	}
}

func (uq *userQuery) Login(username string) (user.Core, error) {
	if username == "" {
		log.Println("data empty, query error")
		return user.Core{}, errors.New("username is empty")
	}

	res := User{}
	if err := uq.db.Where("username = ?", username).First(&res).Error; err != nil {
		log.Println("login query error", err.Error())
		return user.Core{}, errors.New("data not found")
	}

	return ToCore(res), nil
}

func (uq *userQuery) Register(newUser user.Core) (user.Core, error) {
	if newUser.Username == "" || newUser.Password == "" {
		log.Println("data empty")
		return user.Core{}, errors.New("username or password is empty")
	}

	dupUser := CoreToData(newUser)
	err := uq.db.Where("username = ?", newUser.Username).First(&dupUser).Error
	if err == nil {
		log.Println("duplicated")
		return user.Core{}, errors.New("username duplicated")
	}

	cnv := CoreToData(newUser)
	err = uq.db.Create(&cnv).Error
	if err != nil {
		log.Println("query error", err.Error())
		return user.Core{}, errors.New("server error")
	}

	newUser.ID = cnv.ID
	return newUser, nil
}

func (uq *userQuery) Profile(userID uint) (user.Core, error) {
	res := User{}
	if err := uq.db.Where("id=?", userID).First(&res).Error; err != nil {
		log.Println("get profile query error", err.Error())
		return user.Core{}, errors.New("query error, problem with server")
	}

	return ToCore(res), nil
}

func (uq *userQuery) GetByEmail(email string) (user.Core, error) {
	model := User{}
	if err := uq.db.Where("email = ?", email).Take(&model).Error; err != nil {
		return user.Core{}, err
	}

	return ToCore(model), nil
}

func (uq *userQuery) Update(userID uint, updateData user.Core) (user.Core, error) {
	cnv := CoreToData(updateData)
	res := User{}
	qry := uq.db.Model(&res).Where("id = ?", userID).Updates(&cnv)

	affrows := qry.RowsAffected
	if affrows == 0 {
		log.Println("no rows affected")
		return user.Core{}, errors.New("no data updated")
	}

	err := qry.Error
	if err != nil {
		log.Println("update user query error", err.Error())
		return user.Core{}, err
	}

	return ToCore(cnv), nil
}

func (uq *userQuery) Delete(userID uint) error {
	res := User{}
	qry := uq.db.Delete(&res, userID)

	rowAffect := qry.RowsAffected
	if rowAffect <= 0 {
		log.Println("no data processed")
		return errors.New("no user has delete")
	}

	err := qry.Error
	if err != nil {
		log.Println("delete query error", err.Error())
		return errors.New("delete account fail")
	}

	return nil
}
