package user

import (
	"campyuk-api/config"
	"campyuk-api/features/user"
	"campyuk-api/pkg/helper"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

type userService struct {
	qry     user.UserRepository
	vld     *validator.Validate
	storage user.StorageGateway
	g       user.GoogleGateway
}

func New(ud user.UserRepository, vld *validator.Validate, storage user.StorageGateway, g user.GoogleGateway) user.UserService {
	return &userService{
		qry:     ud,
		vld:     vld,
		storage: storage,
		g:       g,
	}
}

func (us *userService) Register(newUser user.Core) (user.Core, error) {
	hashed, err := helper.GeneratePassword(newUser.Password)
	if err != nil {
		log.Println("bcrypt error ", err.Error())
		return user.Core{}, errors.New("password process error")
	}

	err = us.vld.Struct(&newUser)
	if err != nil {
		log.Println("err", err)
		msg := helper.ValidationErrorHandle(err)
		return user.Core{}, errors.New(msg)
	}

	if newUser.Role == "admin" {
		return user.Core{}, errors.New("cannot register as admin")
	}

	newUser.Password = string(hashed)
	res, err := us.qry.Register(newUser)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "duplicated") {
			msg = "data already used or duplicated"
		} else if strings.Contains(err.Error(), "empty") {
			msg = "username not allowed empty"
		} else {
			msg = "server error"
		}
		return user.Core{}, errors.New(msg)
	}

	return res, nil
}

func (us *userService) Login(username, password string) (string, user.Core, error) {
	res, err := us.qry.Login(username)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "empty") {
			msg = "username or password not allowed empty"
		} else {
			msg = "account not registered or server error"
		}
		return "", user.Core{}, errors.New(msg)
	}

	if err := helper.CheckPassword(res.Password, password); err != nil {
		log.Println("login compare", err.Error())
		return "", user.Core{}, errors.New("password not matched")
	}

	useToken, _ := helper.GenerateJWT(int(res.ID), res.Role)
	return useToken, res, nil
}

func (us *userService) LoginGoogle(accessToken string, refreshToken string) (user.Core, error) {
	email, err := us.g.GetEmail(accessToken)
	if err != nil {
		log.Println(err)
		return user.Core{}, errors.New("internal server error")
	}

	core, err := us.qry.GetByEmail(email)
	if err != nil {
		log.Println(err)
		msg := "internal server error"
		if strings.Contains(err.Error(), "not found") {
			msg = "user not found"
		}
		return user.Core{}, errors.New(msg)
	}

	if core.Role == "admin" {
		// Store token in local storage
		f, err := os.Create(config.TokenPath)
		if err != nil {
			log.Println(err)
			return user.Core{}, errors.New("internal server errro")
		}
		defer f.Close()
		token := make(map[string]interface{})
		token["refresh_token"] = refreshToken
		data, err := json.Marshal(token)
		if err != nil {
			log.Println(err)
			return user.Core{}, errors.New("internal server errro")
		}
		_, err = f.Write(data)
		if err != nil {
			log.Println(err)
			return user.Core{}, errors.New("internal server errro")
		}
	}

	return core, nil
}

func (us *userService) Profile(token interface{}) (user.Core, error) {
	id, _ := helper.ExtractToken(token)
	if id <= 0 {
		return user.Core{}, errors.New("data not found")
	}

	res, err := us.qry.Profile(id)
	if err != nil {
		log.Println("data not found")
		return user.Core{}, errors.New("query error, problem with server")
	}

	return res, nil
}

func (us *userService) Update(token interface{}, fileData *multipart.FileHeader, updateData user.Core) (user.Core, error) {
	id, _ := helper.ExtractToken(token)
	if updateData.Password != "" {
		hashed, _ := helper.GeneratePassword(updateData.Password)
		updateData.Password = string(hashed)
	}

	if fileData != nil {
		secureURL, err := us.storage.Upload(fileData)
		if err != nil {
			log.Println(err)
			var msg string
			if strings.Contains(err.Error(), "bad request") {
				msg = err.Error()
			} else {
				msg = "failed to upload image, server error"
			}
			return user.Core{}, errors.New(msg)
		}
		updateData.UserImage = secureURL
	}

	res, err := us.qry.Update(uint(id), updateData)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server error"
		}
		return user.Core{}, errors.New(msg)
	}

	return res, nil
}

func (us *userService) Delete(token interface{}) error {
	id, _ := helper.ExtractToken(token)
	err := us.qry.Delete(uint(id))

	if err != nil {
		log.Println("query error", err.Error())
		return errors.New("query error, delete account fail")
	}

	return nil
}
