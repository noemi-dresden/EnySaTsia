package main

import (
	"fmt"
	"meMySelf/EnySaTsia/database"
	"meMySelf/EnySaTsia/models"
	"testing"

	"github.com/kataras/iris/httptest"
)

//TestNewApp tests the routes in the application
func TestApp(t *testing.T) {
	app := App()
	e := httptest.New(t, app)

	// test if app is reachable
	e.GET("/").Expect().Status(httptest.StatusOK)

	// register user without role
	userWithoutRoleAndEmail := models.User{}
	userWithoutRoleAndEmail.Email = ""
	userWithoutRoleAndEmail.Password = "za"
	userWithoutRoleAndEmail.Role = ""

	e.POST("/user/register").
		WithJSON(userWithoutRoleAndEmail).
		Expect().
		Status(httptest.StatusNotAcceptable).
		Body().
		Contains("user must have a role and email")

	fmt.Print("\nSTORY: Register \n AS A unregistered user\n Scenario 1:\n GIVEN a user that wants to register\n" +
		"WHEN the Role or(and) email are missing\n THEN the service returns a StatusNotAcceptable\n")

	// register user
	user := models.User{}
	user.Email = "za"
	user.Password = "za"
	user.Role = "za"
	e.POST("/user/register").
		WithJSON(user).Expect().Status(httptest.StatusOK).Body().Contains("success")

	fmt.Print("\nSTORY: Register \n AS A unregistered user\n Scenario 2:\n GIVEN a user that wants to register\n" +
		"WHEN the Role or(and) email are not missing\n THEN the service returns a StatusOK and a message 'success' \n")

	// register existing user
	e.POST("/user/register").
		WithJSON(user).
		Expect().
		Status(httptest.StatusConflict).
		Body().
		Contains("User exists already")

	fmt.Print("\nSTORY: Register \n AS A unregistered user\n Scenario 3:\n GIVEN a user that user wants to register\n" +
		"WHEN the Email is already registered\n THEN the service returns a StatusConflict and a message 'User exists already'\n")

	// login user that is in registered
	u := models.User{}
	u.Email = "za"
	u.Password = "za"
	e.POST("/user/login").
		WithJSON(u).
		Expect().
		Status(httptest.StatusOK).
		Body().
		Contains("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiemEifQ.o-kLje8VzI4Wq7oLXT-ytWmDMAuP64H-bO5b_TaAO2s")

	fmt.Print("\nSTORY: Login \n AS A unlogged user\n Scenario 1:\n GIVEN a user that wants to login\n" +
		"WHEN the user is registered and Password matchs\n THEN the service returns a StatusOK and a message cintaining the jwt Token\n")

	// login user that is not registered
	us := models.User{}
	us.Email = "zo"
	us.Password = "za"
	e.POST("/user/login").
		WithJSON(us).
		Expect().
		Status(httptest.StatusUnauthorized).
		Body().
		Contains("No user")

	fmt.Print("\nSTORY: Login \n AS A unregistered user\n Scenario 1:\n GIVEN a user that wants to login\n" +
		"WHEN the user is not registered\n THEN the service returns a StatusUnauthorized and a message 'No User'\n")

	// login user with false password
	usWithFalsePassword := models.User{}
	usWithFalsePassword.Email = "za"
	usWithFalsePassword.Password = "zo"
	e.POST("/user/login").
		WithJSON(usWithFalsePassword).
		Expect().
		Status(httptest.StatusUnauthorized).
		Body().
		Contains("Password not match")

	fmt.Print("\nSTORY: Login \n AS A unlogged user\n Scenario 2:\n GIVEN a user that wants to login\n" +
		"WHEN the password does not match\n THEN the service returns a StatusUnauthorized and a message 'Password not match'\n")

	// change password without Authorization
	newPassword := models.PasswordChange{}
	newPassword.Email = "za"
	newPassword.OldPassword = "zo"
	newPassword.NewPassword = "zo"
	e.POST("/user/changePassword").Expect().Status(httptest.StatusUnauthorized)

	fmt.Print("\nSTORY: Change password \n AS A user\n Scenario 1:\n GIVEN a user that wants to change password\n" +
		"WHEN the user has no Authorizationd\n THEN the service returns a StatusUnauthorized\n")

	// change password with wrong Authorization
	e.POST("/user/changePassword").
		WithJSON(newPassword).
		WithHeader("Authorization", "Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiemEifQ.o-kLje8VzI4Wq7oLXT-ytWmDMAuP64H-bO5b_TaAO2s").
		Expect().
		Status(httptest.StatusUnauthorized)

	fmt.Print("\nSTORY: Change password \n AS A user\n Scenario 2:\n GIVEN a user that wants to change password\n" +
		"WHEN the user has wrong Authorization \n THEN the service returns a StatusUnauthorized\n")

	// change password with wrong old password
	e.POST("/user/changePassword").
		WithJSON(newPassword).
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiemEifQ.o-kLje8VzI4Wq7oLXT-ytWmDMAuP64H-bO5b_TaAO2s").
		Expect().
		Status(httptest.StatusForbidden).
		Body().
		Contains("Not your old password")

	fmt.Print("\nSTORY: Change password \n AS A user\n Scenario 3:\n GIVEN a user that wants to change password\n" +
		"WHEN the old password does not match \n THEN the service returns a StatusForbidden and a message 'Not your old password'\n")

	// change password of non existing user
	newPassword.Email = "zo"
	e.POST("/user/changePassword").
		WithJSON(newPassword).
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiemEifQ.o-kLje8VzI4Wq7oLXT-ytWmDMAuP64H-bO5b_TaAO2s").
		Expect().
		Status(httptest.StatusNotFound).
		Body().
		Contains("User not found")

	fmt.Print("\nSTORY: Change password \n AS A user\n Scenario 4:\n GIVEN a user that wants to change password\n" +
		"WHEN the user does not exists \n THEN the service returns a StatusNotFound and a message 'User not found'\n")

	// change password with right Authorization
	newPassword.Email = "za"
	newPassword.OldPassword = "za"
	e.POST("/user/changePassword").
		WithJSON(newPassword).
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiemEifQ.o-kLje8VzI4Wq7oLXT-ytWmDMAuP64H-bO5b_TaAO2s").
		Expect().
		Status(httptest.StatusOK).
		Body().
		Contains("success")
	fmt.Print("\nSTORY: Change password \n AS A user\n Scenario 5:\n GIVEN a user that wants to change password\n" +
		"WHEN the old password does match \n AND the user has an Authorization \n AND the Authorization is not wrong \n AND user exists \n THEN the service returns a StatusOK and a message 'Success'\n")

	// remove user from database
	err := database.RemoveUser("za")
	if err == nil {
		fmt.Print("Test user deleted")
	}
}
