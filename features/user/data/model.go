package data

import (
	"campyuk-api/features/camp/data"
	"campyuk-api/features/user"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	Fullname  string
	Email     string
	UserImage string
	Password  string
	Role      string
	Camps     []data.Camp `gorm:"foreignKey:HostID"`
}

func ToCore(data User) user.Core {
	return user.Core{
		ID:        data.ID,
		Username:  data.Username,
		Fullname:  data.Fullname,
		Email:     data.Email,
		UserImage: data.UserImage,
		Password:  data.Password,
		Role:      data.Role,
	}
}

func CoreToData(data user.Core) User {
	return User{
		Model:     gorm.Model{ID: data.ID},
		Username:  data.Username,
		Fullname:  data.Fullname,
		Email:     data.Email,
		UserImage: data.UserImage,
		Password:  data.Password,
		Role:      data.Role,
	}
}
