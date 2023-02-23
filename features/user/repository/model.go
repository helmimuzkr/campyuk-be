package repository

import (
	booking "campyuk-api/features/booking/repository"
	camp "campyuk-api/features/camp/repository"
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
	Camps     []camp.Camp       `gorm:"foreignKey:HostID"`
	Bookings  []booking.Booking `gorm:"foreignKey:UserID"`
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
