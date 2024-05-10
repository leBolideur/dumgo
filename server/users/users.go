package users

import (
	"crypto/md5"
	"dumgo/server/db"
	"dumgo/server/utils"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Token    []byte
	ReqCount int
	Db       *db.DumDB
	Nick     string
}

type UserPool struct {
	Users []*User
}

func (up *UserPool) isNickExists(nick string) bool {
	for _, user := range up.Users {
		if user.Nick == nick {
			return true
		}
	}

	return false
}

type UserChanRespData struct {
	Success bool
	User    *User
	Msg     string
}

var UserChan = make(chan string)
var UserChanResp = make(chan UserChanRespData)

type UserComm int

func (uc *UserComm) LogInUser(args *UserCommArgs, reply *UserCommResponse) error {
	user := AddUser(args.Nick)
	msg := fmt.Sprintf("User %s created with token %s", user.Nick, user.Token)
	*reply = UserCommResponse{Status: true, Token: msg}
	return nil
}

func AddUser(nick string) *User {
	UserChan <- fmt.Sprintf("log %s", nick)
	resp := <-UserChanResp
	return resp.User
}

func HandleUser(pool *UserPool, ch chan string, chResp chan UserChanRespData) {
	for {
		msg := <-ch
		split := strings.Split(msg, " ")
		switch split[0] {
		case "log":
			nick := split[1]
			h := md5.New()
			io.WriteString(h, nick)
			utils.Logger.Println("User log request for ", nick)

			isNickExists := pool.isNickExists(nick)
			if isNickExists {
				utils.Logger.Println("Nickname already in use: ", nick)
				ch <- ""
			}

			newUser := &User{
				Token:    (h.Sum(nil)),
				ReqCount: 0,
				Db:       nil,
				Nick:     nick,
			}

			pool.Users = append(pool.Users, newUser)
			chResp <- UserChanRespData{Success: true, User: newUser, Msg: "User logged in"}
		}
	}
}
