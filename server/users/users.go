package users

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"dumgo/server/db"
	"dumgo/server/utils"
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

func (up *UserPool) isLogged(token string) *User {
	for _, user := range up.Users {
		if string(user.Token) == token {
			return user
		}
	}

	return nil
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
	resp := AddUser(args.Nick)
	if !resp.Success {
		msg := fmt.Errorf("User auth failed > %s", resp.Msg)
		*reply = UserCommResponse{Status: false, Token: "", Msg: msg.Error()}
		return msg
	}

	user := resp.User
	msg := fmt.Sprintf("User %s logged in with token %s", user.Nick, user.Token)
	*reply = UserCommResponse{Status: true, Token: string(user.Token), Msg: msg}
	return nil
}

func AddUser(nick string) UserChanRespData {
	UserChan <- fmt.Sprintf("log %s", nick)
	resp := <-UserChanResp
	return resp
}

func HandleUser(pool *UserPool, ch chan string, chResp chan UserChanRespData) {
	for {
		msg := <-ch
		split := strings.Split(msg, " ")
		switch split[0] {
		case "get":
			token := split[1]
			user := pool.isLogged(token)

			var resp UserChanRespData
			if user == nil {
				resp = UserChanRespData{Success: false, User: nil, Msg: "Not logged"}
			} else {
				resp = UserChanRespData{Success: true, User: user, Msg: "Ok"}
			}

			chResp <- resp
		case "log":
			nick := split[1]
			h := md5.New()
			io.WriteString(h, nick)
			utils.Logger.Println("User log request for ", nick)

			isNickExists := pool.isNickExists(nick)
			if isNickExists {
				chResp <- UserChanRespData{Success: false, User: nil, Msg: "Nickname already in use"}
			}

			newUser := &User{
				Token:    h.Sum(nil)[:2],
				ReqCount: 0,
				Db:       new(db.DumDB),
				Nick:     nick,
			}

			pool.Users = append(pool.Users, newUser)
			chResp <- UserChanRespData{Success: true, User: newUser, Msg: "User logged in"}
		}
	}
}
