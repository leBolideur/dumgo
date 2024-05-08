package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"dumgo/pkgs/db"
)

func main() {
	logger := log.New(os.Stdout, "[dumgo]", log.Ldate|log.Lshortfile)

	dumdb := db.NewDumDB()
	rpc.Register(dumdb)

	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("listener error: ", err)
		return
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		go rpc.ServeConn(con)
		logger.Println("request received!")
	}
}
