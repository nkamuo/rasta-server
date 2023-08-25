package auth

import (
	"errors"

	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/utils/token"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/argon2"
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(user model.User, password string) (authToken string, err error) {

	userRepo := repository.GetUserRepository()

	userPassword := user.Password

	if nil == userPassword {
		userPassword, err = userRepo.GetPassword(&user)
		if err != nil {
			return "", err
		}
	}

	err = VerifyPassword(password, userPassword.HashedPassword)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// return "", err
		return "", errors.New("username or password not correct")
	}

	authToken, err = token.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return authToken, nil

}
