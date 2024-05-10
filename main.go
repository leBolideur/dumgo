package main

import (
	"net"
	"net/rpc"

	"dumgo/server/request"
	"dumgo/server/users"
	"dumgo/server/utils"
)

func main() {
	userComm := new(users.UserComm)
	rpc.Register(userComm)
	rpc.Register(new(request.Request))

	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		utils.Logger.Println("listener error: ", err)
		return
	}
	defer listener.Close()

	userPool := users.UserPool{
		Users: []*users.User{},
	}

	go users.HandleUser(&userPool, users.UserChan, users.UserChanResp)
	request.HandleRequest(listener)
}
