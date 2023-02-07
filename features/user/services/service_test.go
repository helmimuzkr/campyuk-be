package services

import (
	"bytes"
	"campyuk-api/features/user"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	data := mocks.NewUserData(t)
	input := user.Core{Username: "griffin", Fullname: "griffinhenry", Email: "grf29@gmail.com", Password: "gg123"}
	resData := user.Core{ID: uint(1), Username: "griffin", Fullname: "griffinhenry", Email: "grf29@gmail.com"}
	srv := New(data)

	t.Run("success create account", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(resData, nil).Once()
		res, err := srv.Register(input)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		assert.NotEmpty(t, resData.Username)
		data.AssertExpectations(t)
	})

	t.Run("username not allowed empty", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.Core{}, errors.New("data not allowed to empty")).Once()
		res, err := srv.Register(input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "empty")
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.Core{}, errors.New("server error")).Once()
		res, err := srv.Register(input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server error")
		data.AssertExpectations(t)
	})

	t.Run("data already used", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.Core{}, errors.New("data already used, duplicated")).Once()
		res, err := srv.Register(input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "already used")
		data.AssertExpectations(t)
	})

	t.Run("bcrypt error", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.Core{}, errors.New("password processed error")).Once()
		res, err := srv.Register(input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "error")
		data.AssertExpectations(t)
	})

	// t.Run("minimal 5 character", func(t *testing.T) {
	// 	inputData := user.Core{Username: "grf", Fullname: "griffinhenry", Email: "grf29@gmail.com", Password: "123"}
	// 	data.On("Register", mock.Anything).Return(user.Core{}, errors.New("input value must be greater")).Once()
	// 	res, err := srv.Register(inputData)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, uint(0), res.ID)
	// 	assert.ErrorContains(t, err, "greater")
	// 	data.AssertExpectations(t)
	// })

	// t.Run("format email", func(t *testing.T) {
	// 	data.On("Register", mock.Anything).Return(user.Core{}, errors.New("input value must be an email")).Once()
	// 	res, err := srv.Register(input)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, uint(0), res.ID)
	// 	assert.ErrorContains(t, err, "email")
	// 	data.AssertExpectations(t)
	// })
}

func TestLogin(t *testing.T) {
	data := mocks.NewUserData(t)
	input := "grf@gmail.com"
	hashed, _ := helper.GeneratePassword("gg123")
	resData := user.Core{ID: uint(1), Username: "griffin", Fullname: "griffinhenry", Email: "grf@gmail.com", Password: hashed}

	t.Run("success login", func(t *testing.T) {
		data.On("Login", input).Return(resData, nil).Once()
		srv := New(data)
		token, res, err := srv.Login(input, "gg123")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, resData.Username, res.Username)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Login", input).Return(user.Core{}, errors.New("server error")).Once()
		srv := New(data)
		_, res, err := srv.Login(input, "gg123")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error")
		assert.Empty(t, nil)
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("username or password empty", func(t *testing.T) {
		data.On("Login", input).Return(user.Core{}, errors.New("username or password not allowed empty")).Once()
		srv := New(data)
		_, res, err := srv.Login(input, "")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "empty")
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("login compare", func(t *testing.T) {
		data.On("Login", input).Return(resData, errors.New("password not matched")).Once()
		srv := New(data)
		token, res, err := srv.Login(input, "123")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not matched")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})
}

func TestProfile(t *testing.T) {
	data := mocks.NewUserData(t)
	resData := user.Core{ID: 1, Username: "griffin", Fullname: "griffinhenry", Email: "grf@gmail.com"}
	srv := New(data)

	t.Run("success show profile", func(t *testing.T) {
		data.On("Profile", uint(1)).Return(resData, nil).Once()

		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Profile(useToken)
		assert.Nil(t, err)
		assert.Equal(t, resData.Fullname, res.Fullname)
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("Profile", uint(1)).Return(user.Core{}, errors.New("query error, problem with server")).Once()

		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Profile(useToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, user.Core{}, res)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Profile", uint(1)).Return(user.Core{}, errors.New("query error, problem with server")).Once()

		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Profile(useToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error")
		assert.Equal(t, user.Core{}, res)
		data.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	data := mocks.NewUserData(t)
	inputData := user.Core{ID: uint(1), Username: "griffin", Fullname: "griffinhenry", Email: "grf@gmail.com"}
	resData := user.Core{ID: uint(1), Username: "griffinh", Fullname: "griffinhenry", Email: "grif@gmail.com"}

	t.Run("success add image", func(t *testing.T) {
		srv := New(data)
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

		_, err = srv.Update(pToken, header, inputData)
		if err != nil {
			log.Println(err.Error())
		}
	})

	t.Run("failed to upload image", func(t *testing.T) {
		srv := New(data)
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

		_, err = srv.Update(pToken, header, inputData)
		if err != nil {
			log.Println(err.Error())
		}

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
	})

	t.Run("success update account", func(t *testing.T) {
		data.On("Update", uint(1), inputData).Return(resData, nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, &multipart.FileHeader{}, inputData)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("Update", uint(1), inputData).Return(user.Core{}, errors.New("not found")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, &multipart.FileHeader{}, inputData)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, user.Core{}, res)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Update", uint(1), inputData).Return(user.Core{}, errors.New("server error")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, &multipart.FileHeader{}, inputData)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, user.Core{}, res)
		data.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	data := mocks.NewUserData(t)

	t.Run("success delete profile", func(t *testing.T) {
		data.On("Delete", uint(1)).Return(nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Delete", mock.Anything).Return(errors.New("server error")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error")
		data.AssertExpectations(t)
	})
}
