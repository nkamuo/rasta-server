package dto

type UserCreationInput struct {
	UserRegistrationInput
	IsAdmin bool `json:"isAdmin" binding:"optional"`
}

type UserUpdateInput struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type UserOutput struct {
	Email   string `json:"email" binding:"required"`
	IsAdmin bool   `json:"isAdmin" binding:"optional"`
	//PROFILE INFORMATION
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}
