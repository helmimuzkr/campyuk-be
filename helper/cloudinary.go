package helper

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"campyuk-api/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Uploader interface {
	Upload(file *multipart.FileHeader) (string, error)
	Destroy(publicID string) error
}

type claudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinary(cfg *config.AppConfig) Uploader {
	cld, err := cloudinary.NewFromParams(cfg.CLOUDINARY_CLOUD_NAME, cfg.CLOUDINARY_API_KEY, cfg.CLOUDINARY_API_SECRET)
	if err != nil {
		log.Println("init cloudinary gagal", err)
		return nil
	}

	return &claudinaryUploader{cld: cld}
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

func (cu *claudinaryUploader) Upload(file *multipart.FileHeader) (string, error) {
	src, _ := file.Open()
	defer src.Close()

	publicID := fmt.Sprintf("%d-%s", int(file.Size), time.Now().Format("20060102-150405")) // Format  "file_size-(YY-MM-DD)-(hh-mm-ss)""

	uploadResult, err := cu.cld.Upload.Upload(
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

func (cu *claudinaryUploader) Destroy(publicID string) error {
	_, err := cu.cld.Upload.Destroy(
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
