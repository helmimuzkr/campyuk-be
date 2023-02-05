package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func FileHeaderTest(t *testing.T) {

	// // Open file needed
	docSrc, err := os.Open("./file-test/document.pdf")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer docSrc.Close()
	imgSrc1, err := os.Open("./file-test/image1.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer imgSrc1.Close()
	imgSrc2, err := os.Open("./file-test/image2.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer imgSrc2.Close()
	imgSrc3, err := os.Open("./file-test/image3.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer imgSrc3.Close()

	// Prepare request file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	// docDst, err := writer.CreateFormFile("document", "./file-test/document.pdf")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	imgDst1, err := writer.CreateFormFile("images", "./file-test/image1.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	imgDst2, err := writer.CreateFormFile("images", "./file-test/image2.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	imgDst3, err := writer.CreateFormFile("images", "./file-test/image3.png")
	if err != nil {
		log.Fatal(err.Error())
	}

	// _, err = io.Copy(docDst, docSrc)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = io.Copy(imgDst1, imgSrc1)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = io.Copy(imgDst2, imgSrc2)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = io.Copy(imgDst3, imgSrc3)
	if err != nil {
		log.Fatal(err.Error())
	}

	req, _ := http.NewRequest("POST", "http://localhost:8000/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	form := req.MultipartForm.File
	documentHeader := form["document"][0]
	imagesHeader := form["images"]

	// _, documentHeader, err := req.FormFile("document")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	fmt.Println(documentHeader.Filename)
	fmt.Println(imagesHeader)

	t.Run("Success add new camp", func(t *testing.T) {
	})

}
