package user

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID        uint
	Username  string `validate:"min=5"`
	Fullname  string `validate:"required"`
	Email     string `validate:"required,email"`
	UserImage string
	Password  string `validate:"min=5"`
	Role      string `validate:"min=4"`
}

type UserHandler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Profile() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Update() echo.HandlerFunc
	GoogleAuth() echo.HandlerFunc
	GoogleCallback() echo.HandlerFunc
}

type UserService interface {
	Login(username, password string) (string, Core, error)
	LoginGoogle(accessToken string, refreshToken string) (Core, error)
	Register(newUser Core) (Core, error)
	Profile(token interface{}) (Core, error)
	Update(token interface{}, fileData *multipart.FileHeader, updateData Core) (Core, error)
	Delete(token interface{}) error
}

type UserRepository interface {
	Login(username string) (Core, error)
	Register(newUser Core) (Core, error)
	Profile(userID uint) (Core, error)
	GetByEmail(email string) (Core, error)
	Update(userID uint, updateData Core) (Core, error)
	Delete(userID uint) error
}

type GoogleGateway interface {
	GetEmail(accessToken string) (string, error)
}

type StorageGateway interface {
	Upload(file *multipart.FileHeader) (string, error)
	Destroy(secureURL string) error
}
