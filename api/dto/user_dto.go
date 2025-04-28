package dto

type UserDto struct {
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type CreateUserDto struct {
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type AssignUserToGroupDto struct {
	UserId  string `json:"userId" binding:"required"`
	GroupId string `json:"groupId" binding:"required"`
}

type AssignUserToGroupsDto struct {
	UserId   string   `json:"userId" binding:"required"`
	GroupIds []string `json:"groupIds" binding:"required"`
}

type UnassignUserFromGroupDto struct {
	UserId  string `json:"userId" binding:"required"`
	GroupId string `json:"groupId" binding:"required"`
}

type UnassignUserFromGroupsDto struct {
	UserId   string   `json:"userId" binding:"required"`
	GroupIds []string `json:"groupIds" binding:"required,min=1"`
}

type LoginDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordDto struct {
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserDto struct {
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type LogoutDto struct {
	RefreshToken struct {
		Token string `json:"token" binding:"required"`
	} `json:"refreshToken" binding:"required"`
}

type RoleMappingPayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
