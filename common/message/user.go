package message

import (
	"net"
)

type User struct {
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` //200在线100离线
}

type UserList map[int]OnlineUser

type OnlineUser struct {
	Conn net.Conn
	User
}
