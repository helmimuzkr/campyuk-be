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

func setupTest(t *testing.T) (*mocks.Uploader, *mocks.CampData, camp.CampService) {
	v := validator.New()
	up := mocks.NewUploader(t)
	data := mocks.NewCampData(t)
	srv := New(data, v, up)

	return up, data, srv
}

func dataSample() camp.Core {
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

	return inData
}

func TestAdd(t *testing.T) {
	upload, data, srv := setupTest(t)
	inData := dataSample()

	t.Run("Success add camp", func(t *testing.T) {
		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()
		upload.On("Upload", &multipart.FileHeader{Filename: "image.jpg"}).Return("www.cloudinary.com/image.jpg", nil).Once()

		inData.Document = "www.cloudinary.com/document.pdf"
		inData.Images = []camp.Image{{ImageURL: "www.cloudinary.com/image.jpg"}}
		data.On("Add", uint(1), inData).Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{{Filename: "image.jpg"}})

		assert.Nil(t, err)
		upload.AssertExpectations(t)
		data.AssertExpectations(t)
	})

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

	t.Run("not pdf", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document"}, []*multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request because of format not pdf")
	})

	t.Run("not image", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{{Filename: "image.svg"}})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request because of format not png, jpg, or jpeg")
	})

	t.Run("failed to upload document", func(t *testing.T) {
		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("", errors.New("failed to upload document because internal server error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{{Filename: "image.jpg"}})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to upload document")
		upload.AssertExpectations(t)
	})

	t.Run("failed to upload image", func(t *testing.T) {
		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()
		upload.On("Upload", &multipart.FileHeader{Filename: "image.jpg"}).Return("", errors.New("failed to upload image because internal server error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{{Filename: "image.jpg"}})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to upload image")
		upload.AssertExpectations(t)
	})

	// t.Run("failed to destroy image", func(t *testing.T) {
	// 	upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()
	// 	upload.On("Upload", &multipart.FileHeader{Filename: "image.jpg"}).Return("", errors.New("failed to upload image because internal server error")).Once()
	// 	upload.On("Destroy", "file/image").Return(errors.New("failed to upload document because internal server error")).Once()

	// 	_, tkn := helper.GenerateJWT(1, "host")
	// 	token := tkn.(*jwt.Token)
	// 	token.Valid = true

	// 	err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{{Filename: "image.jpg"}})

	// 	assert.NotNil(t, err)
	// 	assert.ErrorContains(t, err, "failed to upload image")
	// 	upload.AssertExpectations(t)
	// })

	t.Run("internal server error", func(t *testing.T) {
		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()
		upload.On("Upload", &multipart.FileHeader{Filename: "image.jpg"}).Return("www.cloudinary.com/image.jpg", nil).Once()

		inData.Document = "www.cloudinary.com/document.pdf"
		inData.Images = []camp.Image{{ImageURL: "www.cloudinary.com/image.jpg"}}
		inData.VerificationStatus = "PENDING"
		data.On("Add", uint(1), inData).Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Add(token, inData, &multipart.FileHeader{Filename: "document.pdf"}, []*multipart.FileHeader{{Filename: "image.jpg"}})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		upload.AssertExpectations(t)
		data.AssertExpectations(t)

	})
}

func TestList(t *testing.T) {
	_, data, srv := setupTest(t)
	inData := dataSample()

	listRes := []camp.Core{inData, inData}
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
	_, data, srv := setupTest(t)
	resData := dataSample()

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
	upload, data, srv := setupTest(t)
	inData := dataSample()

	resData := inData
	resData.Items = []camp.CampItem{{ID: uint(1), Name: "Tenda", Stock: 5, RentPrice: 10000}}

	t.Run("Success update camp", func(t *testing.T) {
		oldData := camp.Core{Document: ""}
		data.On("GetByID", uint(1), uint(1)).Return(oldData, nil).Once()

		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()

		inData.Document = "www.cloudinary.com/document.pdf"
		data.On("Update", uint(1), uint(1), inData).Return(nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.Nil(t, err)
		upload.AssertExpectations(t)
		data.AssertExpectations(t)
	})

	t.Run("access denied", func(t *testing.T) {
		_, tkn := helper.GenerateJWT(1, "guest")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "access is denied due to invalid credential")
	})

	t.Run("not pdf", func(t *testing.T) {
		oldData := camp.Core{Document: ""}
		data.On("GetByID", uint(1), uint(1)).Return(oldData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "image.jpg"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request because of format not pdf")
	})

	t.Run("Get by id not found", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(camp.Core{}, errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "camp not found")
		data.AssertExpectations(t)
	})

	t.Run("Get by id database error", func(t *testing.T) {
		data.On("GetByID", uint(1), uint(1)).Return(camp.Core{}, errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})

	t.Run("failed to upload document", func(t *testing.T) {
		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("", errors.New("failed to upload document because internal server error")).Once()

		oldData := camp.Core{Document: "document.pdf"}
		data.On("GetByID", uint(1), uint(1)).Return(oldData, nil).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to upload document")
		upload.AssertExpectations(t)
	})

	t.Run("Update camp not found", func(t *testing.T) {
		oldData := camp.Core{Document: "old-document.pdf"}
		data.On("GetByID", uint(1), uint(1)).Return(oldData, nil).Once()

		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()

		inData.Document = "www.cloudinary.com/document.pdf"
		data.On("Update", uint(1), uint(1), inData).Return(errors.New("not found")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "camp not found")
		upload.AssertExpectations(t)
	})

	t.Run("Update database error", func(t *testing.T) {
		oldData := camp.Core{Document: "old-document.pdf"}
		data.On("GetByID", uint(1), uint(1)).Return(oldData, nil).Once()

		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()

		inData.Document = "www.cloudinary.com/document.pdf"
		data.On("Update", uint(1), uint(1), inData).Return(errors.New("query error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		upload.AssertExpectations(t)
	})

	t.Run("Destroy old image in cloudinary error", func(t *testing.T) {
		oldData := camp.Core{Document: "cloudinary/file/old-document.pdf"}
		data.On("GetByID", uint(1), uint(1)).Return(oldData, nil).Once()

		upload.On("Upload", &multipart.FileHeader{Filename: "document.pdf"}).Return("www.cloudinary.com/document.pdf", nil).Once()

		inData.Document = "www.cloudinary.com/document.pdf"
		data.On("Update", uint(1), uint(1), inData).Return(nil).Once()

		upload.On("Destroy", "file/old-document").Return(errors.New("failed to upload document because internal server error")).Once()

		_, tkn := helper.GenerateJWT(1, "host")
		token := tkn.(*jwt.Token)
		token.Valid = true

		err := srv.Update(token, uint(1), inData, &multipart.FileHeader{Filename: "document.pdf"})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to destroy")
		upload.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	_, data, srv := setupTest(t)

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
	_, data, srv := setupTest(t)

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
