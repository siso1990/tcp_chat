package process

import (
	"encoding/json"
	"fmt"
	"net"

	"go_code/chat2.0/chatServer2.0/model"
	"go_code/chat2.0/common/message"
	"go_code/chat2.0/common/utils"
)

type UserProcess struct {
	Conn net.Conn
}

var OnlineUsers message.UserList = make(map[int]message.OnlineUser, 10)

//登录
func (u *UserProcess) SaveProcessLogin(mes *message.Message) (err error) {
	var longinMes message.LoginMessage
	var loginResMes message.LoginMessageRes
	err = json.Unmarshal([]byte(mes.Data), &longinMes)
	if err != nil {
		fmt.Println("登录消息解析出错", err)
		return
	}

	user, err := model.MyUserDao.Login(longinMes.UserId, longinMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_PWD {
			loginResMes.Error = err.Error()
			loginResMes.Code = 400
		} else if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Error = err.Error()
			loginResMes.Code = 500
		}
	} else {
		loginResMes.Code = 200
		fmt.Println(user.UserId, "已登录")
		var onlineUser = message.OnlineUser{u.Conn, *user}
		OnlineUsers[longinMes.UserId] = onlineUser
		fmt.Println(OnlineUsers)
	}
	//通过各种消息获得mes
	mesres := message.MarshalMessage(loginResMes, message.LoginMessageResType)
	err = utils.WritMessage(u.Conn, &mesres)
	if err != nil {
		fmt.Println("发送登录返回消息失败")
		return
	}
	//推送上线消息
	if loginResMes.Code == 200 {
		u.UserStatusProcess(*user, 200)
	}
	return
}

//注册
func (u *UserProcess) SaveProcessRegister(mes *message.Message) (err error) {
	var regmes message.RegisterMessage
	var regmesRes message.RegisterMessageRes
	err = json.Unmarshal([]byte(mes.Data), &regmes)
	if err != nil {
		fmt.Println("解析用户注册信息失败", err)
		return
	}

	err = model.MyUserDao.Register(regmes.User.UserId, regmes.User.UserPwd, regmes.User.UserName)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			regmesRes.Code = 500
			regmesRes.Error = err.Error()
		} else {
			regmesRes.Code = 400
			regmesRes.Error = "系统错误"
			return
		}
	} else {
		regmesRes.Code = 200
	}
	mesres := message.MarshalMessage(regmesRes, message.RegisterMessageResType)
	err = utils.WritMessage(u.Conn, &mesres)
	if err != nil {
		fmt.Println("返回用户注册信息失败", err)
		return
	}
	return

}

//在线用户列表
func (u *UserProcess) SaveProcessUserList(mes *message.Message) (err error) {
	var userListmes message.UserListMessage
	var userListmesRes message.UserListMessageRes
	err = json.Unmarshal([]byte(mes.Data), &userListmes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&userListmes)", err)
		return
	}
	userListmesRes.UserList = make(map[int]message.User, 10)
	for i, v := range OnlineUsers {
		userListmesRes.UserList[i] = v.User
	}

	fmt.Println(userListmesRes.UserList)
	if len(OnlineUsers) != 0 {
		userListmesRes.Code = 200
	} else {
		userListmesRes.Code = 500
	}
	mesres := message.MarshalMessage(userListmesRes, message.UserListMessageResType)
	err = utils.WritMessage(u.Conn, &mesres)
	return
}

//上线通知
func (u *UserProcess) UserStatusProcess(user message.User, status int) (err error) {
	var userStatus message.UserStatus
	userStatus.User = user
	userStatus.Status = status
	mes := message.MarshalMessage(userStatus, message.UserStatusType)
	for _, v := range OnlineUsers {
		utils.WritMessage(v.Conn, &mes)
	}
	return err
}
