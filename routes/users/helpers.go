package routes

import (
	"EnySaTsia/database"
	"EnySaTsia/models"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt - Hash and Salt password
func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Print(err)
	}
	return string(hash)
}

// ComparePassword - user Passwort vs entered password by the user
func ComparePassword(password string, compareTo string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(compareTo))
	if err != nil {
		fmt.Print(err)
		return false
	}
	return true
}

// CheckUser - Check if the email is already registeres
func CheckUser(email string) (models.User, bool) {
	return database.FindUser(email)
}

// CreateToken - Create Token for signed user
func CreateToken(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": user.Role,
	})
	tokenString, _ := token.SignedString(key)
	return tokenString
}
