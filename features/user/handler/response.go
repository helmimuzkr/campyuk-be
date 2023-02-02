package handler

import "campyuk-api/features/user"

type UserReponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	UserImage string `json:"user_image"`
	Role      string `json:"role"`
}

func ToResponse(data user.Core) UserReponse {
	return UserReponse{
		ID:       data.ID,
		Username: data.Username,
		Fullname: data.Fullname,
		Email:    data.Email,
		Role:     data.Role,
	}
}

func GetToResponse(data user.Core) UserReponse {
	return UserReponse{
		ID:        data.ID,
		Username:  data.Username,
		Fullname:  data.Fullname,
		Email:     data.Email,
		UserImage: data.UserImage,
		Role:      data.Role,
	}
}
