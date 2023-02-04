package service

import (
	"campyuk-api/features/camp"
	"campyuk-api/helper"
	"errors"
	"log"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
)

type campService struct {
	qry camp.CampData
	vld *validator.Validate
}

func New(q camp.CampData, v *validator.Validate) camp.CampService {
	return &campService{
		qry: q,
		vld: v,
	}
}

func (cs *campService) Add(token interface{}, newCamp camp.Core, document *multipart.FileHeader, imagesHeader []*multipart.FileHeader) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	if err := cs.vld.Struct(newCamp); err != nil {
		msg := helper.ValidationErrorHandle(err)
		return errors.New(msg)
	}

	docURL, err := helper.UploadFile(document)
	if err != nil {
		log.Println(err)
		var msg string
		if strings.Contains(err.Error(), "bad request") {
			msg = err.Error()
		} else {
			msg = "failed to upload image because internal server error"
		}
		return errors.New(msg)
	}

	imageCore := []camp.Image{}
	for _, h := range imagesHeader {
		imageURL, err := helper.UploadFile(h)
		if err != nil {
			log.Println(err)
			var msg string
			if strings.Contains(err.Error(), "bad request") {
				msg = err.Error()
			} else {
				msg = "failed to upload image because internal server error"
			}

			// Hapus image di Cloudinary(terlanjur upload) jika salah satu image gagal diupload
			for _, v := range imageCore {
				publicID := helper.GetPublicID(v.ImageURL)
				if err = helper.DestroyFile(publicID); err != nil {
					log.Println(err)
					return errors.New("failed to upload image because internal server error")
				}
			}
			return errors.New(msg)
		}

		imageCore = append(imageCore, camp.Image{ImageURL: imageURL})
	}

	newCamp.Document = docURL
	newCamp.Images = imageCore
	newCamp.VerificationStatus = "PENDING"
	if err := cs.qry.Add(userID, newCamp); err != nil {
		return errors.New("internal server error")
	}

	return nil
}

func (cs *campService) List(token interface{}) ([]camp.Core, error) {
	id, role := helper.ExtractToken(token)

	res, err := cs.qry.List(id, role)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return res, nil
}

func (cs *campService) GetByID(token interface{}, campID uint) (camp.Core, error) {
	userID, role := helper.ExtractToken(token)

	res, err := cs.qry.GetByID(userID, campID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "camp not found"
		} else {
			msg = "internal server errorr"
		}
		return camp.Core{}, errors.New(msg)
	}

	if role != "host" && role != "admin" {
		res.Document = ""
	}

	return res, nil
}

func (cs *campService) Update(token interface{}, campID uint, updateCamp camp.Core, document *multipart.FileHeader) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	res, err := cs.qry.GetByID(userID, campID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "camp not found"
		} else {
			msg = "internal server errorr"
		}
		return errors.New(msg)
	}

	if document != nil {
		docURL, err := helper.UploadFile(document)
		if err != nil {
			log.Println(err)
			var msg string
			if strings.Contains(err.Error(), "bad request") {
				msg = err.Error()
			} else {
				msg = "failed to upload image because internal server error"
			}
			return errors.New(msg)
		}
		updateCamp.Document = docURL
	}

	if err := cs.qry.Update(userID, campID, updateCamp); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "camp not found"
		} else {
			msg = "internal server errorr"
		}
		return errors.New(msg)
	}

	if res.Document != "" {
		publicID := helper.GetPublicID(res.Document)
		if err := helper.DestroyFile(publicID); err != nil {
			log.Println("destroy file", err)
			return errors.New("failed to destroy image")
		}
	}

	return nil
}

func (cs *campService) Delete(token interface{}, campID uint) error {
	userID, _ := helper.ExtractToken(token)

	err := cs.qry.Delete(userID, campID)
	if err != nil {
		log.Println("delete error")
		if strings.Contains(err.Error(), "cannot") {
			return errors.New("access is denied")
		}
		return errors.New("internal server error")
	}
	return nil
}

func (cs *campService) RequestAdmin(token interface{}, campID uint) error {
	return nil
}
