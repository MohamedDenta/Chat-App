package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	connectedUsers     map[int]*User
	Messages           []*Message `json:messages`
	addUser            chan *User
	removeUser         chan *User
	newIncomingMessage chan *Message
	errorChannel       chan error
	doneCh             chan bool
}

func NewServer() *Server {
	Messages := []*Message{}
	connectedUsers := make(map[int]*User)
	addUser := make(chan *User)
	removeUser := make(chan *User)
	newIncomingMessage := make(chan *Message)
	errorChannel := make(chan error)
	doneCh := make(chan bool)

	return &Server{
		connectedUsers:     connectedUsers,
		Messages:           Messages,
		addUser:            addUser,
		removeUser:         removeUser,
		newIncomingMessage: newIncomingMessage,
		errorChannel:       errorChannel,
		doneCh:             doneCh,
	}
}

func (server *Server) AddUser(user *User) {
	log.Println("In AddUser")
	server.addUser <- user
}

func (server *Server) RemoveUser(user *User) {
	log.Println("Removing user")
	server.removeUser <- user
}

func (server *Server) ProcessNewIncomingMessage(message *Message) {
	log.Println("In ProcessNewIncomingMessage ", message)
	server.newIncomingMessage <- message
}

func (server *Server) Done() {
	server.doneCh <- true
}

func (server *Server) sendPastMessages(user *User) {
	for _, msg := range server.Messages {
		log.Println("In sendPastMessages writing ", msg)
		user.Write(msg)
	}
}

func (server *Server) Err(err error) {
	server.errorChannel <- err
}

func (server *Server) sendAll(msg *Message) {
	log.Println("In Sending to all Connected users")
	for _, user := range server.connectedUsers {
		user.Write(msg)
	}
}

func (server *Server) Listen() {
	log.Println("Server Listening ... ")
	http.HandleFunc("/chat", server.handleChat)
	http.HandleFunc("/getAllMessages", server.handleGetAllMessages)

	for {
		select {
		case user := <-server.addUser:
			log.Println("added a new User")
			server.connectedUsers[user.id] = user
			log.Println("Now ", len(server.connectedUsers), " users are connected to chat room")
			server.sendPastMessages(user)
		case user := <-server.removeUser:
			log.Println("Removing user from chat room")
			delete(server.connectedUsers, user.id)
		case msg := <-server.newIncomingMessage:
			server.Messages = append(server.Messages, msg)
			server.sendAll(msg) // -- incase of chat room
		case err := <-server.errorChannel:
			log.Println("Error : ", err)
		case <-server.doneCh: // to stop server
			return
		}
	}
}

func (server *Server) handleChat(responseWrite http.ResponseWriter, request *http.Request) {
	log.Println("Handling chat request ")
	var messageObject Message
	conn, _ := upgrader.Upgrade(responseWrite, request, nil)

	err := conn.ReadJSON(&messageObject)
	log.Println("Message retrieved when add user recieved", &messageObject)

	if err != nil {
		log.Println("Error while reading JSON from websocket ", err.Error())
	}
	user := NewUser(conn, server)

	log.Println("going to add user", user)
	server.AddUser(user)

	log.Println("user added successfully")
	server.ProcessNewIncomingMessage(&messageObject)
	user.Listen()

}

func (server *Server) handleGetAllMessages(responseWriter http.ResponseWriter, request *http.Request) {

	json.NewEncoder(responseWriter).Encode(server)
}
