package database

import (
	"EnySaTsia/models"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//var urlDb = "localhost:27017" "mongodb://mongo:27017"
var urlDb = "10.100.30.78:27017"

//Database name
var Database = "voting-app"
var mgoSession *mgo.Session

//Connect --connect to database
func Connect() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(urlDb)
		if err != nil {
			log.Fatal("Failed to start mongo db")
		}
	}
	return mgoSession.Clone()
}

// FindUser --find a user given email in the database
func FindUser(email string) (models.User, bool) {
	session := Connect()
	defer session.Close()
	c := session.DB(Database).C("users")
	result := models.User{}
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return models.User{}, false
	}
	return result, true
}

// InsertUser -- insert user in the database
func InsertUser(user models.User) (err error) {
	session := Connect()
	defer session.Close()
	c := session.DB(Database).C("users")
	return c.Insert(user)
}

// UpdateUser -- update the user in database
func UpdateUser(user models.User) (err error) {
	session := Connect()
	defer session.Close()
	c := session.DB(Database).C("users")
	return c.UpdateId(user.ID, user)
}

// Remove user with mail
func RemoveUser(email string) (err error) {
	session := Connect()
	defer session.Close()
	c := session.DB(Database).C("users")
	return c.Remove(bson.M{"email": email})
}
