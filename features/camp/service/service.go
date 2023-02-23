package service

import (
	"campyuk-api/features/camp"
	"campyuk-api/pkg/helper"
	"errors"
	"log"
	"math"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
)

type campService struct {
	qry     camp.CampRepository
	vld     *validator.Validate
	storage camp.StorageGateway
}

func New(q camp.CampRepository, v *validator.Validate, storage camp.StorageGateway) camp.CampService {
	return &campService{
		qry:     q,
		vld:     v,
		storage: storage,
	}
}

func (cs *campService) Add(token interface{}, newCamp camp.Core, document *multipart.FileHeader, imagesHeader []*multipart.FileHeader) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	if err := cs.vld.Struct(newCamp); err != nil {
		log.Println(err)
		msg := helper.ValidationErrorHandle(err)
		return errors.New(msg)
	}

	filedoc := strings.Split(document.Filename, ".")
	format := filedoc[len(filedoc)-1]
	if format != "pdf" {
		return errors.New("bad request because of format not pdf")
	}

	for _, img := range imagesHeader {
		fileimg := strings.Split(img.Filename, ".")
		format := fileimg[len(fileimg)-1]
		if format != "png" && format != "jpg" && format != "jpeg" {
			return errors.New("bad request because of format not png, jpg, or jpeg")
		}
	}

	docURL, err := cs.storage.Upload(document)
	if err != nil {
		log.Println(err)
		return errors.New("failed to upload document because internal server error")
	}

	imageCore := []camp.Image{}
	for _, h := range imagesHeader {
		imageURL, err := cs.storage.Upload(h)
		if err != nil {
			log.Println(err)
			return errors.New("failed to upload image because internal server error")
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

func (cs *campService) List(token interface{}, page int) (map[string]interface{}, []camp.Core, error) {
	id, role := helper.ExtractToken(token)

	if page < 1 {
		page = 1
	}
	limit := 4
	// Calculate offset
	offset := (page - 1) * limit

	// Get total record, list camp, and error
	totalRecord, res, err := cs.qry.List(id, role, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal server error")
	}

	totalPage := int(math.Ceil(float64(totalRecord) / float64(limit)))
	if page > totalPage {
		return nil, nil, errors.New("page not found")
	}

	pagination := make(map[string]interface{})
	pagination["page"] = page
	pagination["limit"] = limit
	pagination["offset"] = offset
	pagination["totalRecord"] = totalRecord
	pagination["totalPage"] = totalPage

	return pagination, res, nil
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
		filedoc := strings.Split(document.Filename, ".")
		format := filedoc[len(filedoc)-1]
		if format != "pdf" {
			return errors.New("bad request because of format not pdf")
		}

		docURL, err := cs.storage.Upload(document)
		if err != nil {
			log.Println(err)
			return errors.New("failed to upload document because internal server error")
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

	if document != nil && res.Document != "" {
		if err := cs.storage.Destroy(res.Document); err != nil {
			log.Println("destroy file", err)
			return errors.New("failed to destroy document")
		}
	}

	return nil
}

func (cs *campService) Delete(token interface{}, campID uint) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	// res, err := cs.qry.GetByID(userID, campID)
	// if err != nil {
	// 	log.Println(err)
	// 	msg := ""
	// 	if strings.Contains(err.Error(), "not found") {
	// 		msg = "camp not found"
	// 	} else {
	// 		msg = "internal server errorr"
	// 	}
	// 	return errors.New(msg)
	// }

	err := cs.qry.Delete(userID, campID)
	if err != nil {
		log.Println("delete error")
		if strings.Contains(err.Error(), "not found") {
			return errors.New("camp not found")
		}
		return errors.New("internal server error")
	}

	// if res.Document != "" {
	// 	publicID := helper.GetPublicID(res.Document)
	// 	if err := helper.DestroyFile(publicID); err != nil {
	// 		log.Println("destroy file", err)
	// 		return errors.New("failed to destroy document")
	// 	}
	// }

	// if res.Images != nil {
	// 	for _, v := range res.Images {
	// 		publicID := helper.GetPublicID(v.ImageURL)
	// 		if err := helper.DestroyFile(publicID); err != nil {
	// 			log.Println("destroy file", err)
	// 			return errors.New("failed to destroy image")
	// 		}
	// 	}
	// }

	return nil
}

func (cs *campService) RequestAdmin(token interface{}, campID uint, status string) error {
	_, role := helper.ExtractToken(token)
	if role != "admin" {
		return errors.New("access is denied due to invalid credential")
	}

	if err := cs.qry.RequestAdmin(campID, status); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "camp not found"
		} else {
			msg = "internal server errorr"
		}
		return errors.New(msg)
	}

	return nil
}
