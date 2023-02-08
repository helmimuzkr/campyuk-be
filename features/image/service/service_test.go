package service

import (
	"campyuk-api/features/image"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"

	"mime/multipart"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	data := mocks.NewImageData(t)
	up := mocks.NewUploader(t)
	srv := New(data, up)
	input := image.Core{
		CampID: uint(1),
		Image:  "www.cloudinary.com/image.jpg",
	}

	t.Run("success add image", func(t *testing.T) {
		up.On("Upload", &multipart.FileHeader{Filename: "image.jpg"}).Return("www.cloudinary.com/image.jpg", nil).Once()

		data.On("Add", uint(1), input).Return(nil).Once()
		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.Add(pToken, uint(1), &multipart.FileHeader{Filename: "image.jpg"})
		assert.Nil(t, err)
		up.AssertExpectations(t)
		data.AssertExpectations(t)
	})

	t.Run("failed to upload image", func(t *testing.T) {
		up.On("Upload", &multipart.FileHeader{Filename: "image.jpg"}).Return("", errors.New("failed to upload image")).Once()

		// data.On("Add", uint(1), input).Return(errors.New("failed to upload image")).Once()
		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.Add(pToken, uint(1), &multipart.FileHeader{Filename: "image.jpg"})
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		up.AssertExpectations(t)
		data.AssertExpectations(t)
	})

	t.Run("format not allowed", func(t *testing.T) {
		up.On("Upload", &multipart.FileHeader{Filename: "image.sh"}).Return("", errors.New("bad request because of format not pdf, png, jpg, or jpeg")).Once()

		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.Add(pToken, uint(1), &multipart.FileHeader{Filename: "image.sh"})
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "format")
		up.AssertExpectations(t)
		data.AssertExpectations(t)
	})

	// t.Run("internal server error", func(t *testing.T) {
	// 	f, err := os.Open("/mnt/c/project/campyuk/docs/erd.jpg")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// 	defer f.Close()

	// 	// prepare request body
	// 	// reserve a form field with 'file' as key
	// 	// then assign the file content to field using 'io.Copy'
	// 	// create a http post request, set content type to multipart-form
	// 	// read the 'file' field using 'req.FormFile'

	// 	body := &bytes.Buffer{}
	// 	writer := multipart.NewWriter(body)
	// 	part, err := writer.CreateFormFile("file", "/mnt/c/project/campyuk/docs/erd.jpg")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	_, err = io.Copy(part, f)
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	writer.Close()

	// 	req, _ := http.NewRequest("POST", "/upload", body)
	// 	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 	_, header, _ := req.FormFile("file")

	// 	_, token := helper.GenerateJWT(1, "host")
	// 	pToken := token.(*jwt.Token)
	// 	pToken.Valid = true

	// 	err = srv.Add(pToken, uint(1), header)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}

	// 	assert.NotNil(t, err)
	// 	assert.ErrorContains(t, err, "server")
	// })

	// t.Run("access is denied", func(t *testing.T) {
	// 	f, err := os.Open("/mnt/c/project/campyuk/docs/erd.jpg")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// 	defer f.Close()

	// 	// prepare request body
	// 	// reserve a form field with 'file' as key
	// 	// then assign the file content to field using 'io.Copy'
	// 	// create a http post request, set content type to multipart-form
	// 	// read the 'file' field using 'req.FormFile'

	// 	body := &bytes.Buffer{}
	// 	writer := multipart.NewWriter(body)
	// 	part, err := writer.CreateFormFile("file", "/mnt/c/project/campyuk/docs/erd.jpg")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	_, err = io.Copy(part, f)
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	writer.Close()

	// 	req, _ := http.NewRequest("POST", "/upload", body)
	// 	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 	_, header, _ := req.FormFile("file")

	// 	_, token := helper.GenerateJWT(1, "user")
	// 	pToken := token.(*jwt.Token)
	// 	pToken.Valid = true

	// 	err = srv.Add(pToken, uint(1), header)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}

	// 	assert.NotNil(t, err)
	// 	assert.ErrorContains(t, err, "denied")
	// })

	// t.Run("format not allowed", func(t *testing.T) {
	// 	f, err := os.Open("/mnt/c/project/campyuk/test.sh")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// 	defer f.Close()

	// 	// prepare request body
	// 	// reserve a form field with 'file' as key
	// 	// then assign the file content to field using 'io.Copy'
	// 	// create a http post request, set content type to multipart-form
	// 	// read the 'file' field using 'req.FormFile'

	// 	body := &bytes.Buffer{}
	// 	writer := multipart.NewWriter(body)
	// 	part, err := writer.CreateFormFile("file", "/mnt/c/project/campyuk/test.sh")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	_, err = io.Copy(part, f)
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	writer.Close()

	// 	req, _ := http.NewRequest("POST", "/upload", body)
	// 	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 	_, header, _ := req.FormFile("file")

	// 	_, token := helper.GenerateJWT(1, "host")
	// 	pToken := token.(*jwt.Token)
	// 	pToken.Valid = true

	// 	err = srv.Add(pToken, uint(1), header)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}

	// 	assert.NotNil(t, err)
	// 	assert.ErrorContains(t, err, "format")
	// })
}

func TestDelete(t *testing.T) {
	data := mocks.NewImageData(t)
	up := mocks.NewUploader(t)
	srv := New(data, up)

	t.Run("success delete image", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()
		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("server error")).Once()
		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error")
		data.AssertExpectations(t)
	})

	t.Run("invalid credential", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("access is denied")).Once()
		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "denied")
		data.AssertExpectations(t)
	})
}
