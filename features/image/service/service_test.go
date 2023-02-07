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
	input := image.Core{
		ID:     1,
		CampID: 1,
		Image:  "https://res.cloudinary.com/djqjmzwsa/image/upload/v1675603226/campyuk/20230205-212016.jpg",
	}
	resData := image.Core{
		ID:     1,
		CampID: 1,
		Image:  "https://res.cloudinary.com/djqjmzwsa/image/upload/v1675603226/campyuk/20230205-212016.jpg",
	}
	srv := New(data)

	t.Run("success add image", func(t *testing.T) {
		data.On("Add", uint(1), input).Return(resData, nil).Once()

		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Add(pToken, uint(1), &multipart.FileHeader{})
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("format not allowed", func(t *testing.T) {
		data.On("Add", uint(1), input).Return(resData, errors.New("bad request")).Once()

		_, token := helper.GenerateJWT(1, "host")
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Add(pToken, uint(1), &multipart.FileHeader{})
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request")
		data.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	data := mocks.NewImageData(t)

	t.Run("success delete image", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("server error")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1, "user")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error")
		data.AssertExpectations(t)
	})
}
