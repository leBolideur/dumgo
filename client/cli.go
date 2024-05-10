package main

import (
	"fmt"
	"net/rpc"
	"os"
)

func logUser(client *rpc.Client, nick string) {
	type UserCommResponse struct {
		Status bool
		Token  string
	}
	var userReply UserCommResponse
	type LogInArgs struct {
		Nick string
	}
	err := client.Call("UserComm.LogInUser", &LogInArgs{Nick: nick}, &userReply)
	if err != nil {
		fmt.Println("RPC call error: ", err)
		return
	}

	fmt.Printf("reply [%t] >> user token > %x\n", userReply.Status, userReply.Token)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: ./client [request]")
		return
	}

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}
	defer client.Close()

	if len(args) == 3 && args[1] == "log" {
		logUser(client, args[2])
		return
	}

	type ReqResponse struct {
		Success bool
		Msg     string
	}
	var reply ReqResponse

	type ReqArgs struct {
		Request string
	}
	err = client.Call("DumDB.Request", &ReqArgs{Request: args[1]}, &reply)
	if err != nil {
		fmt.Println("RPC call error: ", err)
		return
	}

	fmt.Printf("reply [%t] >> %s\n", reply.Success, reply.Msg)
}
