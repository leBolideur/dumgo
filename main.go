package main

import (
	"net"
	"net/rpc"

	"dumgo/server/db"
	"dumgo/server/users"
	"dumgo/server/utils"
)

func main() {
	dumdb := db.NewDumDB()
	userComm := new(users.UserComm)
	rpc.Register(dumdb)
	rpc.Register(userComm)

	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		utils.Logger.Println("listener error: ", err)
		return
	}
	defer listener.Close()

	userPool := users.UserPool{
		Users: []users.User{},
	}

	go users.HandleUser(&userPool, users.UserChan)

	for {
		con, err := listener.Accept()
		if err != nil {
			utils.Logger.Fatal("accept error")
			continue
		}

		go rpc.ServeConn(con)
		utils.Logger.Println("request received!")
	}
}
