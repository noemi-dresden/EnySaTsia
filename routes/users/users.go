package routes

import (
	"EnySaTsia/database"
	"EnySaTsia/models"
	"fmt"

	"github.com/kataras/iris"
)

var key = []byte("your-secret-key")

// Register a user
func Register(ctx iris.Context) {

	user := models.User{}
	err := ctx.ReadJSON(&user)

	if err != nil {
		fmt.Println("Error on reading form: " + err.Error())
		return
	}
	userResult, isThere := CheckUser(user.Email)

	if !isThere {
		hashed := HashAndSalt(user.Password)
		user.Password = hashed
		if user.Role != "" && user.Email != "" {
			err := database.InsertUser(user)
			if err != nil {
				ctx.JSON("failed to register")
			} else {
				ctx.StatusCode(200)
				ctx.JSON("success")
			}
		} else {
			ctx.StatusCode(406)
			ctx.JSON("user must have a role and email")
		}
	} else {
		fmt.Println(userResult)
		ctx.StatusCode(409)
		ctx.JSON("User exists already")
	}
}

// Login into the system and get some token
func Login(ctx iris.Context) {
	user := models.User{}
	err := ctx.ReadJSON(&user)

	if err != nil {
		ctx.JSON("error while extracting data from client")
	}

	result, isThere := CheckUser(user.Email)

	if isThere {
		if ComparePassword(result.Password, user.Password) {
			token := CreateToken(result)
			ctx.JSON(token)
		} else {
			ctx.StatusCode(401)
			ctx.JSON("Password not match")
		}
	} else {
		ctx.StatusCode(401)
		ctx.JSON("No user")
	}
}

//ChangePassword -Change password of user
func ChangePassword(ctx iris.Context) {
	newPassword := models.PasswordChange{}
	err := ctx.ReadJSON(&newPassword)

	if err != nil {
		ctx.JSON("error while extracting data from client")
	}

	user, isThere := CheckUser(newPassword.Email)

	if isThere {
		match := ComparePassword(user.Password, newPassword.OldPassword)
		if match {
			user.Password = HashAndSalt(newPassword.NewPassword)
			error := database.UpdateUser(user)
			if error != nil {
				ctx.JSON("failed to update")
			} else {
				ctx.StatusCode(200)
				ctx.JSON("success")
			}
		} else {
			ctx.StatusCode(403)
			ctx.JSON("Not your old password")
		}
	} else {
		ctx.StatusCode(404)
		ctx.JSON("User not found")
	}
}
