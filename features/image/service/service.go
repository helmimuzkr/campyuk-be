package service

import (
	"campyuk-api/features/image"
	"campyuk-api/helper"
	"errors"
	"log"
	"mime/multipart"
	"strings"
)

type imageService struct {
	qry image.ImageData
}

func New(repo image.ImageData) image.ImageService {
	return &imageService{qry: repo}
}

func (is *imageService) Add(token interface{}, campID uint, header *multipart.FileHeader) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	imageURL, err := helper.UploadFile(header)
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

	newImage := image.Core{CampID: campID, Image: imageURL}
	if err := is.qry.Add(userID, newImage); err != nil {
		log.Println(err)
		var msg string
		if strings.Contains(err.Error(), "access is denied") {
			msg = err.Error()
		} else {
			msg = "failed to upload image because internal server error"
		}
		return errors.New(msg)
	}

	return nil
}

func (is *imageService) Delete(token interface{}, imageID uint) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	if err := is.qry.Delete(userID, imageID); err != nil {
		log.Println(err)
		var msg string
		if strings.Contains(err.Error(), "access is denied") {
			msg = err.Error()
		} else {
			msg = "internal server error"
		}
		return errors.New(msg)
	}

	return nil
}
