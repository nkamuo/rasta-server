package dto

type UserRegistrationInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	//PROFILE INFORMATION
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type UserCreationInput struct {
	UserRegistrationInput
	IsAdmin bool `json:"isAdmin" binding:"optional"`
}

type UserUpdateInput struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type UserFormLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
