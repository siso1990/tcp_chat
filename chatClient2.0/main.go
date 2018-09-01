package main

import (
	"fmt"
	"os"

	"go_code/chat2.0/chatClient2.0/process"
)

var (
	userId   int
	userPwd  string
	userName string
)

func main() {
	var key int
	for {
		fmt.Println("------多人聊天系统------")
		fmt.Println("------1.登录聊天室------")
		fmt.Println("------2.注册用户------")
		fmt.Println("------3.退出系统------")
		fmt.Printf("请输入选择:")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("------请输入账号------")
			fmt.Scanln(&userId)
			fmt.Println("------请输入密码------")
			fmt.Scanln(&userPwd)
			var up process.UserProcess
			err := up.Login(userId, userPwd)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			fmt.Println("------请输入账号------")
			fmt.Scanln(&userId)
			fmt.Println("------请输入密码------")
			fmt.Scanln(&userPwd)
			fmt.Println("------请输入用户名------")
			fmt.Scanln(&userName)
			var up process.UserProcess
			err := up.Register(userId, userPwd, userName)
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			os.Exit(0)
		default:
			fmt.Println("请输入正确的选项")
		}
	}
}
