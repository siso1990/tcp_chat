// message
package message

import (
	"encoding/json"
	"fmt"
)

type MessageType int

const (
	LoginMessageType       = "LoginMessage"
	LoginMessageResType    = "LoginResMessage"
	RegisterMessageType    = "RegisterMessage"
	RegisterMessageResType = "RegisterMessageRes"
	UserListMessageType    = "UserListMessage"
	UserListMessageResType = "UserListMessageRes"
	UserStatusType         = "UserStatus"
	SmsMessageType         = "SmsMessage"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMessage struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginMessageRes struct {
	Code  int    `json:"code"` //200登录成功,400密码错误,500账户不存在
	Error string `json:"error"`
}

type RegisterMessage struct {
	User User `json:"user"`
}

type RegisterMessageRes struct {
	Code  int    `json:"code"` //200注册成功,500账户已存在
	Error string `json:"error"`
}

type UserListMessage struct {
	UserId int `json:"userId"`
}

type UserListMessageRes struct {
	Code     int          `json:"code"`
	UserList map[int]User `json:"userList"`
}

type UserStatus struct {
	User   User `json:"userId"`
	Status int  `json:"Status"`
}

type SmsMessage struct {
	UserId    int    `json:"userId"`
	Content   string `json:"content"`
	Addressee int    `json:"addressee"` //0代表群发,其余代表各个用户的Id
	Error     string `json:"error"`
}

//各种mes序列化
func MarshalMessage(a interface{}, mesType string) (mes Message) {
	data1, err := json.Marshal(a)
	if err != nil {
		fmt.Println("json.Marshal(a)", err)
		return
	}
	mes.Data = string(data1)
	mes.Type = mesType
	return
}

