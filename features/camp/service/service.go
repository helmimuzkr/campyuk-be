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

	imageURLs := []string{}
	for _, h := range imagesHeader {
		image, err := helper.UploadFile(h)
		if err != nil {
			log.Println(err)
			var msg string
			if strings.Contains(err.Error(), "bad request") {
				msg = err.Error()
			} else {
				msg = "failed to upload image because internal server error"
			}

			// Hapus image di Cloudinary(terlanjur upload) jika salah satu image gagal diupload
			for _, url := range imageURLs {
				publicID := helper.GetPublicID(url)
				if err = helper.DestroyFile(publicID); err != nil {
					log.Println(err)
					return errors.New("failed to upload image because internal server error")
				}
			}
			return errors.New(msg)
		}
		imageURLs = append(imageURLs, image)
	}

	newCamp.Document = docURL
	newCamp.Images = imageURLs
	newCamp.VerificationStatus = "PENDING"
	if err := cs.qry.Add(userID, newCamp); err != nil {
		return errors.New("internal server error")
	}

	return nil
}

func (cs *campService) List(token interface{}) ([]camp.Core, error) {
	return nil, nil
}

func (cs *campService) GetByID(token interface{}, campID uint) (camp.Core, error) {
	return camp.Core{}, nil
}

func (cs *campService) Update(token interface{}, campID uint, udpateCamp camp.Core, document *multipart.FileHeader, image []*multipart.FileHeader) error {
	return nil
}

func (cs *campService) Delete(token interface{}, campID uint) error {
	return nil
}

func (cs *campService) RequestAdmin(token interface{}, campID uint) error {
	return nil
}
