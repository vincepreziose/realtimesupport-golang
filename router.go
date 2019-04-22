package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handler func
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Router struct
type Router struct {
	rules map[string]Handler
}

// NewRouter func
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
	}
}

func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	client := NewClient(socket, r.FindHandler)
	go client.Write()
	client.Read()
}
