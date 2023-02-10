package service

import (
	"campyuk-api/features/item"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	data := mocks.NewItemData(t)
	input := item.Core{ID: 1, Name: "bonfire", Stock: 5, Price: 10000}
	resData := item.Core{ID: 1, Name: "bonfire", Stock: 5, Price: 10000}

	t.Run("sukses tambah data", func(t *testing.T) {
		data.On("Add", uint(1), uint(1), input).Return(resData, nil).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Add(useToken, uint(1), input)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		data.AssertExpectations(t)
	})

	t.Run("invalid credential", func(t *testing.T) {
		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "guest")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Add(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Empty(t, res)
		assert.ErrorContains(t, err, "access is denied")
	})

	t.Run("validation error", func(t *testing.T) {
		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		input := item.Core{}
		res, err := srv.Add(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Empty(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "is required")
	})

	t.Run("item not found", func(t *testing.T) {
		data.On("Add", uint(1), uint(1), input).Return(item.Core{}, errors.New("data not found")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Add(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "not found")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Add", uint(1), uint(1), input).Return(item.Core{}, errors.New("server error")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Add(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		data.AssertExpectations(t)
	})

	t.Run("access is denied", func(t *testing.T) {
		data.On("Add", uint(1), uint(1), input).Return(item.Core{}, errors.New("access denied")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Add(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "denied")
		data.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	data := mocks.NewItemData(t)
	input := item.Core{ID: uint(1), Name: "bonfire", Stock: 5, Price: 10000}
	resData := item.Core{ID: uint(1), Name: "sleepingbag", Stock: 10, Price: 20000}

	t.Run("success update item", func(t *testing.T) {
		data.On("Update", uint(1), uint(1), input).Return(resData, nil).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Update(useToken, uint(1), input)
		assert.Nil(t, err)
		assert.NotEqual(t, input.Name, res.Name)
		assert.NotEqual(t, input.Price, res.Price)
		data.AssertExpectations(t)
	})

	t.Run("invalid credential", func(t *testing.T) {
		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "guest")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Update(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied")
		assert.Empty(t, res)
	})

	t.Run("item not found", func(t *testing.T) {
		data.On("Update", uint(1), uint(1), input).Return(item.Core{}, errors.New("data not found")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Update(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "not found")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Update", uint(1), uint(1), input).Return(item.Core{}, errors.New("server error")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Update(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		data.AssertExpectations(t)
	})

	t.Run("access is denied", func(t *testing.T) {
		data.On("Update", uint(1), uint(1), input).Return(item.Core{}, errors.New("access denied")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		res, err := srv.Update(useToken, uint(1), input)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "denied")
		data.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	data := mocks.NewItemData(t)

	t.Run("success delete item", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("invalid credential", func(t *testing.T) {
		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "guest")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied")
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("Delete", uint(5), uint(1)).Return(errors.New("data not found")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(5, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		data.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("terdapat masalah pada server")).Once()
		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		data.AssertExpectations(t)
	})

	t.Run("access is denied", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("access denied")).Once()

		v := validator.New()
		srv := New(data, v)

		_, token := helper.GenerateJWT(1, "host")
		useToken := token.(*jwt.Token)
		useToken.Valid = true
		err := srv.Delete(useToken, uint(1))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "denied")
		data.AssertExpectations(t)
	})
}
