package db

import (
	"encoding/json"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

var DB *leveldb.DB
var Open bool

type UserDB struct {
	Id   string `json:"email"`
	Name string `json:"username"`
}
type MessageDB struct {
	UserName  string `json:"userName"` // Email
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

func opendatabase() bool {
	if !Open {
		Open = true
		DBpath := "Database/users"
		var err error
		DB, err = leveldb.OpenFile(DBpath, nil)
		if err != nil {
			return false
		}
		return true
	}
	return true
}

// AddUser insert user to database
func AddUser(user UserDB) bool {
	if !Open {
		opendatabase()
	}
	d, _ := json.Marshal(user)
	err := DB.Put([]byte(user.Id), d, nil)
	if err != nil {
		log.Println("failed to add user to database ", err.Error())
		return false
	}
	log.Println("added user to database ")
	return true
}

// AddMsg add message to database
func AddMsg(msg MessageDB) bool {
	if !Open {
		opendatabase()
	}
	d, _ := json.Marshal(msg)
	err := DB.Put([]byte(msg.UserName), d, nil)
	if err != nil {
		log.Println("failed to add msg to database ", err.Error())
		return false
	}
	log.Println("added msg to database ")
	return true
}

// GetAllMsg get all msgs from database
func GetAllMsg() (msgs []MessageDB) {
	if !Open {
		opendatabase()
	}
	iter := DB.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		var newdata MessageDB
		json.Unmarshal(value, &newdata)
		msgs = append(msgs, newdata)
	}
	return msgs
}
