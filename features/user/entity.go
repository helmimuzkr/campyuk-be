package user

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID        uint
	Username  string
	Fullname  string
	Email     string
	UserImage string
	Password  string
	Role      string
}

type UserHandler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Profile() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Update() echo.HandlerFunc
}

type UserService interface {
	Login(username, password string) (string, Core, error)
	Register(newUser Core) (Core, error)
	Profile(token interface{}) (Core, error)
	Update(token interface{}, fileData multipart.FileHeader, updateData Core) (Core, error)
	Delete(token interface{}) error
}

type UserData interface {
	Login(username string) (Core, error)
	Register(newUser Core) (Core, error)
	Profile(userID uint) (Core, error)
	Update(id uint, updateData Core) (Core, error)
	Delete(id uint) error
}
