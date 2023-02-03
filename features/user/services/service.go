package services

import (
	"campyuk-api/features/user"
	"campyuk-api/helper"
	"errors"
	"log"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
)

type userUseCase struct {
	qry user.UserData
	vld *validator.Validate
}

func New(ud user.UserData) user.UserService {
	return &userUseCase{
		qry: ud,
		vld: validator.New(),
	}
}

func (uuc *userUseCase) Register(newUser user.Core) (user.Core, error) {
	hashed, err := helper.GeneratePassword(newUser.Password)
	if err != nil {
		log.Println("bcrypt error ", err.Error())
		return user.Core{}, errors.New("password process error")
	}

	err = uuc.vld.Struct(&newUser)
	if err != nil {
		log.Println("err", err)
		return user.Core{}, errors.New("bad request")
	}

	newUser.Password = string(hashed)
	res, err := uuc.qry.Register(newUser)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "duplicated") {
			msg = "data already used"
		} else if strings.Contains(err.Error(), "empty") {
			msg = "username not allowed empty"
		} else {
			msg = "server error"
		}
		return user.Core{}, errors.New(msg)
	}

	return res, nil
}

func (uuc *userUseCase) Login(email, password string) (string, user.Core, error) {
	res, err := uuc.qry.Login(email)
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

func (uuc *userUseCase) Profile(token interface{}) (user.Core, error) {
	id, _ := helper.ExtractToken(token)
	if id <= 0 {
		return user.Core{}, errors.New("data not found")
	}

	res, err := uuc.qry.Profile(uint(id))
	if err != nil {
		log.Println("data not found")
		return user.Core{}, errors.New("query error, problem with server")
	}

	return res, nil
}

func (uuc *userUseCase) Update(token interface{}, fileData multipart.FileHeader, updateData user.Core) (user.Core, error) {
	id, _ := helper.ExtractToken(token)
	if updateData.Password != "" {
		hashed, _ := helper.GeneratePassword(updateData.Password)
		updateData.Password = string(hashed)
	}
	if fileData.Size != 0 {
		if fileData.Size > 5000000 {
			return user.Core{}, errors.New("size error")
		}
		secureURL, err := helper.UploadFile(&fileData)
		if err != nil {
			log.Println(err)
			var msg string
			if strings.Contains(err.Error(), "wrong input") {
				msg = err.Error()
			} else {
				msg = "failed to upload image, server error"
			}
			return user.Core{}, errors.New(msg)
		}
		updateData.UserImage = secureURL
	}

	res, err := uuc.qry.Update(uint(id), updateData)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server error"
		}
		return user.Core{}, errors.New(msg)
	}

	if res.UserImage != "" {
		publicID := helper.GetPublicID(res.UserImage)
		if err := helper.DestroyFile(publicID); err != nil {
			log.Println("destroy file", err)
			return user.Core{}, errors.New("failed to destroy image")
		}
	}

	return res, nil
}

func (uuc *userUseCase) Delete(token interface{}) error {
	id, _ := helper.ExtractToken(token)
	err := uuc.qry.Delete(uint(id))

	if err != nil {
		log.Println("query error", err.Error())
		return errors.New("query error, delete account fail")
	}

	return nil
}
