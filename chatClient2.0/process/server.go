package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"go_code/chat2.0/common/message"
	"go_code/chat2.0/common/utils"
)

var UserList map[int]message.User = make(map[int]message.User)

var UserInfo message.OnlineUser

func ProcessMes(conn net.Conn) {
	defer conn.Close()
	for {
		mesRes, err := utils.ReadMessage(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch mesRes.Type {
		case message.UserListMessageResType:
			var userListRes message.UserListMessageRes
			err = json.Unmarshal([]byte(mesRes.Data), &userListRes)
			UserList = userListRes.UserList
			fmt.Println("在线用户列表")
			for _, v := range UserList {
				fmt.Println(v.UserId)
			}
		case message.UserStatusType:
			var userStatus message.UserStatus
			err = json.Unmarshal([]byte(mesRes.Data), &userStatus)
			UserList[userStatus.User.UserId] = userStatus.User
			if userStatus.Status == 200 {
				fmt.Println("刚刚上线用户:", userStatus.User.UserId)
			} else if userStatus.Status == 100 {
				fmt.Println("刚刚离线用户:", userStatus.User.UserId)
			}
		case message.SmsMessageType:
			var sms message.SmsMessage
			err = json.Unmarshal([]byte(mesRes.Data), &sms)
			if err != nil {
				fmt.Println("json.Unmarshal([]byte(mesRes.Data), &sms)", err)
				return
			}
			fmt.Println(sms)
			fmt.Println(string(sms.Content), "发送者:", sms.UserId)
		default:
			fmt.Println("未知的消息类型")
		}
	}
}

func Memu() {
	go ProcessMes(UserInfo.Conn)
	for {
		fmt.Printf("--------%v用户已登录-------\n", UserInfo.UserId)
		fmt.Println("--------1.在线用户列表-------")
		fmt.Println("--------2.发送消息-------")
		fmt.Println("--------3.发送文件-------")
		fmt.Println("--------4.退出系统-------")
		var key int
		fmt.Scanln(&key)
		switch key {
		case 1:
			err := UserListMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
		case 2:
			err := SendMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
		case 3:
			
		case 4:
			os.Exit(0)
		default:
			fmt.Println("请输入正确的选项")
		}
	}
}
