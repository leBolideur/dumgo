package request

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"strings"

	"dumgo/server/db"
	"dumgo/server/users"
	"dumgo/server/utils"
)

func HandleRequest(listener net.Listener) {
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

type Request struct{}

func (re *Request) Req(req *db.ReqArgs, reply *db.ReqResponse) error {
	// base := db.NewDumDB()
	cmd := strings.Split(req.Request, " ")
	token := cmd[0]
	if token == "" {
		msg := fmt.Errorf("No token provided")
		*reply = db.ReqResponse{Success: false, Msg: msg.Error()}
		return msg
	}

	users.UserChan <- fmt.Sprintf("get %s", token)
	resp := <-users.UserChanResp
	base := resp.User.Db

	switch cmd[1] {
	case "HEALTH":
		base.Health(reply)
	case "EXPORT":
		base.Export(reply)
	case "RESTORE":
		base.Restore(reply)
	case "SET":
		setArgs := &db.SetArgs{Key: cmd[1], Value: cmd[2]}
		base.Set(setArgs, reply)
	case "GET":
		getArgs := &db.GetArgs{Key: cmd[1]}
		base.Get(getArgs, reply)
	case "INCR":
		base.UpdateInt(cmd[1], "+", 1, reply)
	case "DECR":
		base.UpdateInt(cmd[1], "-", 1, reply)
	case "INCRBY":
		by, err := strconv.ParseInt(cmd[2], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid increment value, got '%s'", cmd[2])
		}
		base.UpdateInt(cmd[1], "+", by, reply)
	case "DECRBY":
		by, err := strconv.ParseInt(cmd[2], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid increment value, got '%s'", cmd[2])
		}
		base.UpdateInt(cmd[1], "-", by, reply)
	default:
		return fmt.Errorf("Unknown cmd '%s'\n", cmd)
	}

	return nil
}
