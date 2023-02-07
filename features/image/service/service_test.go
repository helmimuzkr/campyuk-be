package service

import (
	"bytes"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"mime/multipart"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	data := mocks.NewImageData(t)
	srv := New(data)

	t.Run("success add image", func(t *testing.T) {
		f, err := os.Open("/mnt/c/project/campyuk/docs/erd.jpg")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f.Close()

		// prepare request body
		// reserve a form field with 'file' as key
		// then assign the file content to field using 'io.Copy'
		// create a http post request, set content type to multipart-form
		// read the 'file' field using 'req.FormFile'

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "/mnt/c/project/campyuk/docs/erd.jpg")
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = io.Copy(part, f)
		if err != nil {
			log.Fatal(err.Error())
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, header, _ := req.FormFile("file")

		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err = srv.Add(pToken, uint(1), header)
		if err != nil {
			log.Println(err.Error())
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		f, err := os.Open("/mnt/c/project/campyuk/docs/erd.jpg")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f.Close()

		// prepare request body
		// reserve a form field with 'file' as key
		// then assign the file content to field using 'io.Copy'
		// create a http post request, set content type to multipart-form
		// read the 'file' field using 'req.FormFile'

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "/mnt/c/project/campyuk/docs/erd.jpg")
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = io.Copy(part, f)
		if err != nil {
			log.Fatal(err.Error())
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, header, _ := req.FormFile("file")

		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err = srv.Add(pToken, uint(1), header)
		if err != nil {
			log.Println(err.Error())
		}

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
	})

	t.Run("format not allowed", func(t *testing.T) {
		f, err := os.Open("/mnt/c/project/campyuk/test.sh")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f.Close()

		// prepare request body
		// reserve a form field with 'file' as key
		// then assign the file content to field using 'io.Copy'
		// create a http post request, set content type to multipart-form
		// read the 'file' field using 'req.FormFile'

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "/mnt/c/project/campyuk/test.sh")
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = io.Copy(part, f)
		if err != nil {
			log.Fatal(err.Error())
		}

		writer.Close()

		req, _ := http.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, header, _ := req.FormFile("file")

		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err = srv.Add(pToken, uint(1), header)
		if err != nil {
			log.Println(err.Error())
		}

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "format")
	})
}

func TestDelete(t *testing.T) {
	data := mocks.NewImageData(t)

	t.Run("success delete image", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("server error")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error")
		data.AssertExpectations(t)
	})
}
