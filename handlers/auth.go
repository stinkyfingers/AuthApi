package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stinkyfingers/AuthApi/auth"
	"github.com/stinkyfingers/AuthApi/middleware"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

var publicKey = []byte("key")

func Authenticate(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if token == nil || !token.Valid || err != nil {
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}
	return
}

func Login(ctx *middleware.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var u auth.User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	err = u.Login()
	if err != nil {
		return nil, err
	}
	return json.Marshal(u)

}

func Logout(ctx *middleware.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var u auth.User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	err = u.Logout()
	if err != nil {
		return nil, err
	}
	return json.Marshal(u)

}

func Create(ctx *middleware.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var u auth.User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	err = u.Create()
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return j, err
}
