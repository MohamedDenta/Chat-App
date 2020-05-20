package chat

import (
	"fmt"
	"log"

	"github.com/MohamedDenta/Chat-App/db"
	"github.com/gorilla/websocket"
)

const channelBufSize = 100

var CurrentUser User

type User struct {
	Id              string
	Name            string
	Con             *websocket.Conn
	Servr           *Server
	OutgoingMessage chan *Message
	DoneCh          chan bool
}

func NewUser(conn *websocket.Conn, server *Server) *User {
	if conn == nil {
		panic("connection can not be nil")
	}

	if server == nil {
		panic(" Server cannot be nil")
	}

	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)
	log.Println("Done creating new User")
	return &User{"1", "", conn, server, ch, doneCh}
}

func (user *User) Conn() *websocket.Conn {
	return user.Con
}

func (user *User) Write(message *Message) {
	select {
	case user.OutgoingMessage <- message:
	default:
		user.Servr.RemoveUser(user)
		err := fmt.Errorf("User %d is disconnected.", user.Id)
		user.Servr.Err(err)
	}
}

func (user *User) Done() {
	user.DoneCh <- true
}

func (user *User) Listen() {
	go user.listenWrite()
	user.listenRead()
}

func (user *User) SaveDB() {

	if b := db.AddUser(db.UserDB{Id: user.Id, Name: user.Name}); !b {
		log.Println("error in save user to db ")
		return
	}
}

func (user *User) listenWrite() {
	log.Println("Listening to write to client")

	for {
		select {
		//send message to user
		case msg := <-user.OutgoingMessage:
			log.Println("send in listenWrite for user :", user.Id, msg)
			user.Con.WriteJSON(&msg)

			// receive done request
		case <-user.DoneCh:
			log.Println("Done Channel for user:")
			user.Servr.RemoveUser(user)
			user.DoneCh <- true
			return
		}
	}
}

func (user *User) listenRead() {
	log.Println("Listening to Read to client")
	for {
		select {
		// receive Done request
		case <-user.DoneCh:
			user.Servr.RemoveUser(user)
			user.DoneCh <- true
			return
			// read data from websocket connection
		default:
			var messageObject Message
			err := user.Con.ReadJSON(&messageObject)

			if err != nil {
				user.DoneCh <- true
				log.Println("Error while reading JSON from websocket ", err.Error())
				user.Servr.Err(err)
			} else {
				user.Servr.ProcessNewIncomingMessage(&messageObject)
			}
		}
	}
}
