package service

import (
	"campyuk-api/features/camp"
	"campyuk-api/helper"
	"campyuk-api/mocks"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func SetupTest(t *testing.T) (camp.Core, *mocks.CampData, camp.CampService) {
	v := validator.New()
	data := mocks.NewCampData(t)
	srv := New(data, v)

	inData := camp.Core{
		ID:                 uint(1),
		VerificationStatus: "PENDING",
		Title:              "Tanakita",
		Price:              100000,
		Description:        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi molestie tempus purus, at tristique justo vehicula id.",
		Latitude:           -6.208987101998694,
		Longitude:          106.79970296358913,
		Address:            "Jl. Spartan No.IV, Gotham city, West Java, 53241 +62 985904",
		City:               "Gotham city",
		Distance:           100,
		Document:           "cloudinary.com/document.pdf",
		Images:             []camp.Image{{ID: uint(1), ImageURL: "cloudinary.com//image1.jpg"}, {ID: uint(2), ImageURL: "cloudinary.com//image1.jpg"}},
	}

	return inData, data, srv
}

func TestAdd(t *testing.T) {
	inData, _, srv := SetupTest(t)

	// Successfully add new camp

	t.Run("access denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{}, []*multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
	})

	t.Run("input invalid required", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		inData := camp.Core{}

		err := srv.Add(token, inData, &multipart.FileHeader{}, []*multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "input value")
	})

	t.Run("failed to upload document", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to upload document")
	})

	t.Run("error format document not pdf", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{}, []*multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request")
	})
}

func TestList(t *testing.T) {
	inputData, data, srv := SetupTest(t)

	listRes := []camp.Core{inputData, inputData}
	listRes[1].ID = uint(2)

	t.Run("Success display list camp", func(t *testing.T) {
		data.On("List", uint(1), "guest", 4, 0).Return(2, listRes, nil).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		pagination, actual, err := srv.List(token, 1)

		assert.Nil(t, err)
		assert.Equal(t, listRes[0].ID, actual[0].ID)
		assert.NotNil(t, pagination)
		data.AssertExpectations(t)
	})

	t.Run("error in database", func(t *testing.T) {
		data.On("List", uint(1), "guest", 4, 0).Return(0, nil, errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		pagination, actual, err := srv.List(token, 1)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, pagination)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	resData, data, srv := SetupTest(t)

	resData.Items = []camp.CampItem{{ID: uint(1), Name: "Tenda", Stock: 5, RentPrice: 10000}}

	t.Run("Success get camp", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(resData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetByID(token, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, resData.ID, actual.ID)
		data.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(camp.Core{}, errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetByID(token, uint(1))

		assert.NotNil(t, err)
		assert.Empty(t, actual)
		assert.ErrorContains(t, err, "camp not found")
		data.AssertExpectations(t)
	})

	t.Run("error in database", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(camp.Core{}, errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetByID(token, uint(1))

		assert.NotNil(t, err)
		assert.Empty(t, actual)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	inData, data, srv := SetupTest(t)

	resData := inData
	resData.Items = []camp.CampItem{{ID: uint(1), Name: "Tenda", Stock: 5, RentPrice: 10000}}

	t.Run("access denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
	})

	t.Run("not found", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(camp.Core{}, errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "camp not found")
		data.AssertExpectations(t)
	})

	t.Run("error in database", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(camp.Core{}, errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})

	t.Run("failed to upload document", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(resData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to upload document")
		data.AssertExpectations(t)
	})

	t.Run("error format document not pdf", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(resData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "steam.exe"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request")
		data.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	_, data, srv := SetupTest(t)

	t.Run("Success delete camp", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("access denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")

	})

	t.Run("not found", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "camp not found")
		data.AssertExpectations(t)
	})

	t.Run("error in database", func(t *testing.T) {
		data.On("Delete", uint(1), uint(1)).Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestRequestAdmin(t *testing.T) {
	_, data, srv := SetupTest(t)

	t.Run("Admin successfully accept camp", func(t *testing.T) {
		data.On("RequestAdmin", uint(1), "ACCEPTED").Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "admin")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.RequestAdmin(token, uint(1), "ACCEPTED")

		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("access denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.RequestAdmin(token, uint(1), "ACCEPTED")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
	})

	t.Run("not found", func(t *testing.T) {
		data.On("RequestAdmin", uint(1), "ACCEPTED").Return(errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "admin")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.RequestAdmin(token, uint(1), "ACCEPTED")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "camp not found")
		data.AssertExpectations(t)
	})

	t.Run("error in database", func(t *testing.T) {
		data.On("RequestAdmin", uint(1), "ACCEPTED").Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "admin")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.RequestAdmin(token, uint(1), "ACCEPTED")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

// func TestDelete(t *testing.T) {
// 	inData, data, srv := SetupTest(t)

// 	resData := inData
// 	resData.Items = []camp.CampItem{{ID: uint(1), Name: "Tenda", Stock: 5, RentPrice: 10000}}

// 	t.Run("Success delete camp", func(t *testing.T) {
// 		data.On("GetByID", uint(1), uint(1)).Return(resData, nil).Once()
// 		data.On("Delete", uint(1), uint(1)).Return(nil).Once()

// 		_, tkn := helper.GenerateJWT(1, "host")
// 		token := tkn.(*jwt.Token)
// 		token.Valid = true

// 		err := srv.Delete(token, uint(1))

// 		assert.Nil(t, err)
// 	})

// }
