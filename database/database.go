package database

import (
	"gopkg.in/mgo.v2"

	"os"
	"time"
)

var (
	MongoSession   *mgo.Session
	DatabaseName   string
	AuthCollection = "users"
)

//Init creates mongo session if it does not already exist
func Init() error {
	var err error
	if MongoSession == nil {
		connectionString := mongoConnectionString()
		DatabaseName = connectionString.Database
		MongoSession, err = mgo.DialWithInfo(connectionString)
	}
	return err
}

//Close closes mongo session
func Close() {
	MongoSession.Close()
}

func mongoConnectionString() *mgo.DialInfo {

	var (
		MongoDBHosts    = os.Getenv("MONGO_URL")
		AuthDatabase    = os.Getenv("MONGO_ADMIN_DB")
		AuthUserName    = os.Getenv("MONGO_USERNAME")
		AuthPassword    = os.Getenv("MONGO_PASSWORD")
		mongoDBDialInfo mgo.DialInfo
	)

	if MongoDBHosts == "" {
		mongoDBDialInfo = mgo.DialInfo{
			Addrs:    []string{"127.0.0.1"},
			Timeout:  60 * time.Second,
			Database: "auth",
			Username: "",
			Password: "",
		}
	} else {
		mongoDBDialInfo = mgo.DialInfo{
			Addrs:    []string{MongoDBHosts},
			Timeout:  60 * time.Second,
			Database: AuthDatabase,
			Username: AuthUserName,
			Password: AuthPassword,
		}
	}
	return &mongoDBDialInfo
}
