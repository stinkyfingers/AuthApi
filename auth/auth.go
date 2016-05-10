package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stinkyfingers/AuthApi/database"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"errors"
	"fmt"
	"time"
)

type User struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username,omitempty" bson:"username,omitempty"`
	Password string        `json:"-" bson:"password,omitempty"`
	JWT      string        `json:"jwt,omitempty" bson:"jwt,omitempty"`
	Sites    []string      `json:"sites,omitempty" bson:"sites,omitempty"`
}

var (
	TOKEN_SIGNING_KEY = []byte("key")
)

func init() {
	database.Init()
}

func (u *User) Login() error {
	var err error
	password := u.Password

	//Auth by password
	users, err := u.Find()
	if err != nil || len(users) < 1 {
		if len(users) < 1 {
			err = errors.New("No users with that username.")
		}
		return err
	}

	*u = users[0]
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}

	//set jwt
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	u.JWT, err = token.SignedString(TOKEN_SIGNING_KEY)
	if err != nil {
		return err
	}

	//update user
	u.Password = ""
	err = u.Update()
	return err
}

func (u *User) Logout() error {
	session := database.MongoSession.Clone()
	return session.DB(database.DatabaseName).C(database.AuthCollection).UpdateId(u.ID, bson.M{"$set": bson.M{"jwt": nil}})
}

//Does not match password
func (u *User) Find() ([]User, error) {
	session := database.MongoSession.Clone()
	var result []User
	query := make(map[string]interface{})
	if u.ID.Hex() != "" {
		query["_id"] = u.ID
	}
	if u.Username != "" {
		query["username"] = u.Username
	}
	if u.JWT != "" {
		query["jwt"] = u.JWT
	}

	err := session.DB(database.DatabaseName).C(database.AuthCollection).Find(query).All(&result)
	return result, err
}

func (u *User) Update() error {
	session := database.MongoSession.Clone()
	query := make(map[string]interface{})
	if u.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		query["password"] = string(pass)
	}
	if u.Username != "" {
		query["username"] = u.Username
	}
	if u.JWT != "" {
		query["jwt"] = u.JWT
	}
	if len(u.Sites) > 0 {
		query["sites"] = u.Sites
	}
	return session.DB(database.DatabaseName).C(database.AuthCollection).UpdateId(u.ID, bson.M{"$set": query})
}

func (u *User) Create() error {
	session := database.MongoSession.Clone()
	u.ID = bson.NewObjectId()
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pass)
	return session.DB(database.DatabaseName).C(database.AuthCollection).Insert(u)
}

//VerifyToken confirms validness of token string
func VerifyToken(tokenString string) (jwt.Token, error) {
	var err error
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(TOKEN_SIGNING_KEY), nil
	})
	if err != nil || !token.Valid {
		token = jwt.New(jwt.SigningMethodHS256) //empty token to return
		return *token, fmt.Errorf("Error: %v", err)
	}
	return *token, nil
}

// Authorize returns user with token if okay, else it errors
func (u *User) Authorize() error {
	if err := database.Init(); err != nil {
		return err
	}
	session := database.MongoSession.Clone()
	query := bson.M{
		"jwt": u.JWT,
	}
	return session.DB(database.DatabaseName).C(database.AuthCollection).Find(query).One(&u)
}
