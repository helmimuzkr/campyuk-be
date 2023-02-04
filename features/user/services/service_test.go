package services

import (
	"campyuk-api/features/user"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"
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

	t.Run("success creat account", func(t *testing.T) {
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

	t.Run("account not registered", func(t *testing.T) {
		data.On("Login", input).Return(user.Core{}, errors.New("data not found")).Once()
		srv := New(data)
		token, res, err := srv.Login(input, "gg123")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not registered")
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

// func TestUpdate(t *testing.T) {
// 	data := mocks.NewUserData(t)
// 	inputData := user.Core{ID: uint(1), Username: "griffin", Email: "grf@gmail.com"}
// 	resData := user.Core{ID: uint(1), Username: "griffinh", Email: "grif@gmail.com"}

// 	t.Run("success updating account", func(t *testing.T) {
// 		data.On("Update", uint(1), inputData).Return(resData, nil).Once()
// 		srv := New(data)
// 		_, token := helper.GenerateJWT(1, "user")
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		file, _ := os.Open("./file-test/test.png")
// 		defer file.Close()
// 		fileHeader := &multipart.FileHeader{
// 			Filename: file.Name(),
// 		}

// 		res, err := srv.Update(pToken, *fileHeader, inputData)
// 		assert.Nil(t, err)
// 		assert.NotEqual(t, resData.ID, res.ID)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("incorrect input", func(t *testing.T) {
// 		data.On("Update", uint(1), inputData).Return(user.Core{}, errors.New("incorrect input from user")).Once()
// 		srv := New(data)
// 		_, token := helper.GenerateJWT(1, "user")
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		file, _ := os.Open("./file-test/test.png")
// 		defer file.Close()
// 		fileHeader := &multipart.FileHeader{
// 			Filename: file.Name(),
// 		}

// 		res, err := srv.Update(pToken, *fileHeader, inputData)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "input")
// 		assert.Equal(t, user.Core{}, res)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("failed upload image", func(t *testing.T) {
// 		data.On("Update", uint(1), inputData).Return(user.Core{}, errors.New("failed to upload image, server error")).Once()
// 		srv := New(data)
// 		_, token := helper.GenerateJWT(1, "user")
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		file, _ := os.Open("./file-test/test.png")
// 		defer file.Close()
// 		fileHeader := &multipart.FileHeader{
// 			Filename: file.Name(),
// 		}

// 		res, err := srv.Update(pToken, *fileHeader, inputData)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "server")
// 		assert.Equal(t, user.Core{}, res)
// 		data.AssertExpectations(t)
// 	})
// }

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
