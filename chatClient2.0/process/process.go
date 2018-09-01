package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"go_code/chat2.0/common/message"
	"go_code/chat2.0/common/utils"
)

type UserProcess struct {
	Conn net.Conn
}

func (up *UserProcess) Login(userId int, userPwd string) (err error) {
	up.Conn, err = net.Dial("tcp", "127.0.0.1:8889")
	defer up.Conn.Close()
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//创建登录结构体
	var logmes message.LoginMessage
	logmes.UserId = userId
	logmes.UserPwd = userPwd
	//	//序列化 logmes
	mes := message.MarshalMessage(logmes, message.LoginMessageType)
	//发送登录mes
	err = utils.WritMessage(up.Conn, &mes)
	if err != nil {
		return
	}
	//接受登录返回mes
	mesRes, err := utils.ReadMessage(up.Conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	//返回信息解析
	var loginRes message.LoginMessageRes
	err = json.Unmarshal([]byte(mesRes.Data), &loginRes)
	if loginRes.Code == 200 {
		fmt.Println("账号登录成功")
		//初始化登录用户的信息
		UserInfo = message.OnlineUser{
			Conn: up.Conn,
			User: message.User{
				UserId:     userId,
				UserStatus: 200,
			},
		}
		Memu()
	} else if loginRes.Code == 500 {
		fmt.Println(loginRes.Error)
		os.Exit(0)
	} else if loginRes.Code == 400 {
		fmt.Println(loginRes.Error)
		os.Exit(0)
	}
	return
}

func (up *UserProcess) Register(userId int, userPwd string,
	userName string) (err error) {
	up.Conn, err = net.Dial("tcp", "127.0.0.1:8889")
	defer up.Conn.Close()
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	//创建登录结构体
	var regmes message.RegisterMessage
	regmes.User.UserId = userId
	regmes.User.UserPwd = userPwd
	regmes.User.UserName = userName
	//序列化
	mes := message.MarshalMessage(regmes, message.RegisterMessageType)
	//发送登录mes

	err = utils.WritMessage(up.Conn, &mes)
	if err != nil {
		return
	}
	//接受返回信息
	mesRes, err := utils.ReadMessage(up.Conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	var regRes message.RegisterMessageRes
	err = json.Unmarshal([]byte(mesRes.Data), &regRes)
	if regRes.Code == 200 {
		fmt.Println("账号已经注册成功,请重新登录")
	} else if regRes.Code == 500 {
		fmt.Println(regRes.Error)
	} else {
		fmt.Println(regRes.Error)
	}
	return
}

func UserListMessage() (err error) {
	var userList message.UserListMessage
	//	var userListRes message.UserListMessageRes
	userList.UserId = UserInfo.UserId
	mes := message.MarshalMessage(userList, message.UserListMessageType)
	err = utils.WritMessage(UserInfo.Conn, &mes)
	return nil
}

func SendMessage() (err error) {
	var smsMessage message.SmsMessage
	//	var userListRes message.UserListMessageRes
	smsMessage.UserId = UserInfo.UserId
	var sender int
	fmt.Println("请输入接收者的id,若值0为群发")
	fmt.Scanln(&sender)
	fmt.Println("请输入发送消息的内容,回车发送")
	var content string
	fmt.Scanln(&content)
	smsMessage.Addressee = sender
	smsMessage.Content = content
	mes := message.MarshalMessage(smsMessage, message.SmsMessageType)
	err = utils.WritMessage(UserInfo.Conn, &mes)
	return nil
}
