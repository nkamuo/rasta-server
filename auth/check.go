package auth

import (
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/utils/token"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/argon2"
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(user model.User, password string) (string, error) {

	var err error

	err = VerifyPassword(password, user.HashedPassword)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}
