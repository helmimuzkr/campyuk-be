package pkg

import (
	"campyuk-api/config"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinary(cfg *config.AppConfig) *cloudinaryClient {
	cld, err := cloudinary.NewFromParams(cfg.CLOUDINARY_CLOUD_NAME, cfg.CLOUDINARY_API_KEY, cfg.CLOUDINARY_API_SECRET)
	if err != nil {
		log.Println("init cloudinary gagal", err)
		return nil
	}

	return &cloudinaryClient{cld: cld}
}

func (cc *cloudinaryClient) Upload(file *multipart.FileHeader) (string, error) {
	src, _ := file.Open()
	defer src.Close()

	publicID := fmt.Sprintf("%d-%s", int(file.Size), time.Now().Format("20060102-150405")) // Format  "file_size-(YY-MM-DD)-(hh-mm-ss)""

	uploadResult, err := cc.cld.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			PublicID: publicID,
			Folder:   "file",
		})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func (cc *cloudinaryClient) Destroy(secureURL string) error {
	// Proses filter Public ID dari SecureURL(avatar).
	urls := strings.Split(secureURL, "/")
	urls = urls[len(urls)-2:]                               // array [file, random_name.extension]
	noExtension := strings.Split(urls[len(urls)-1], ".")[0] // remove ".extension", result "random_name"
	urls = append(urls[:1], noExtension)                    // new array [file, random_name]
	publicID := strings.Join(urls, "/")                     // "file/random_name"

	_, err := cc.cld.Upload.Destroy(
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
