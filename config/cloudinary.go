package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(c AppConfig) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET)
	if err != nil {
		log.Println("init cloudinary gagal", err)
		return nil
	}

	return cld
}
