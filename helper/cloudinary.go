package helper

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"campyuk-api/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func NewCloudinary() *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(config.CLOUDINARY_CLOUD_NAME, config.CLOUDINARY_API_KEY, config.CLOUDINARY_API_SECRET)
	if err != nil {
		log.Println("init cloudinary gagal", err)
		return nil
	}

	return cld
}

func UploadFile(file *multipart.FileHeader) (string, error) {
	// Format check
	filename := strings.Split(file.Filename, ".")
	format := filename[len(filename)-1]
	if format != "pdf" && format != "png" && format != "jpg" && format != "jpeg" {
		return "", errors.New("bad request because of format not pdf, png, jpg, or jpeg")
	}

	src, _ := file.Open()
	defer src.Close()

	publicID := time.Now().Format("20060102-150405") // Format  "(YY-MM-DD)-(hh-mm-ss)""

	cld := NewCloudinary()
	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			PublicID: publicID,
			Folder:   config.CLOUDINARY_UPLOAD_FOLDER,
		})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func GetPublicID(secureURL string) string {
	// Proses filter Public ID dari SecureURL(avatar)
	urls := strings.Split(secureURL, "/")
	urls = urls[len(urls)-2:]                               // array [file, random_name.extension]
	noExtension := strings.Split(urls[len(urls)-1], ".")[0] // remove ".extension", result "random_name"
	urls = append(urls[:1], noExtension)                    // new array [file, random_name]
	publicID := strings.Join(urls, "/")                     // "file/random_name"

	return publicID
}

func DestroyFile(publicID string) error {
	// Proses destroy file
	cld := NewCloudinary()
	_, err := cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
