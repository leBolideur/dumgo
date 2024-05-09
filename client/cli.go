package main

import (
	"dumgo/server/db"
	"fmt"
	"net/rpc"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: ./client [request]")
		return
	}

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}
	defer client.Close()

	var reply db.Response

	err = client.Call("DumDB.Request", db.ReqArgs{Request: args[1]}, &reply)
	if err != nil {
		fmt.Println("RPC call error: ", err)
		return
	}

	fmt.Printf("reply [%t] >> %s\n", reply.Success, reply.Msg)
}
