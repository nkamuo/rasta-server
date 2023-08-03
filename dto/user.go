package dto

type UserCreationInput struct {
	UserRegistrationInput
	IsAdmin bool `json:"isAdmin" binding:""`
}

type UserUpdateInput struct {
	Email    *string `json:"email" binding:""`
	Password *string `json:"password" binding:""`
	//PROFILE INFORMATION
	FirstName *string `json:"firstName" binding:""`
	LastName  *string `json:"lastName" binding:""`
	Phone     *string `json:"phone" binding:""`
	IsAdmin   *bool   `json:"isAdmin" binding:""`
}

type UserOutput struct {
	Email   string `json:"email" binding:"required"`
	IsAdmin bool   `json:"isAdmin" binding:""`
	//PROFILE INFORMATION
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}
