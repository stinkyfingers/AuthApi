package middleware

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/stinkyfingers/AuthApi/auth"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

const (
	con = "CONTEXT"
)

type Handler struct {
	Handle      func(*Context, http.ResponseWriter, *http.Request) (interface{}, error)
	Middlewares []Middleware
}

type Middleware struct {
	Handler http.Handler
}
type Context struct {
	MongoSession *mgo.Session
	Params       httprouter.Params
	Data         Data
}
type Data struct {
	Token  string
	UserID bson.ObjectId
}

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}
	ctx := &Context{
		Params: ps,
	}
	context.Set(r, con, ctx)
	for _, middleware := range fn.Middlewares {
		rec := httptest.NewRecorder()
		middleware.Handler.ServeHTTP(rec, r)
		if rec.Code != 200 {
			err := fmt.Errorf("%s", rec.Body.String())
			w.WriteHeader(rec.Code)
			w.Write([]byte(err.Error()))
			return
		}
	}
	ctx = context.Get(r, con).(*Context)

	object, err := fn.Handle(ctx, w, r)

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(object)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	context.Clear(r)
}

//MIDDLEWARE
//Validate By Token
type TokenValidation struct {
	http.Handler
}

func (t TokenValidation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var u auth.User
	u.JWT = r.Header.Get("JWT")
	if u.JWT == "" {
		http.Error(w, "No token", http.StatusInternalServerError)
		return
	}
	u.JWT = strings.TrimPrefix(strings.TrimSpace(u.JWT), "Bearer: ")
	_, err := auth.VerifyToken(u.JWT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = u.Authorize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
