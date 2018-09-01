package main

import (
	"fmt"
	"net"

	"go_code/chat2.0/chatServer2.0/process"
	"go_code/chat2.0/common/message"
	"go_code/chat2.0/common/utils"
)

type Processor struct {
	Conn net.Conn
}

//接受消息,分类处理
func (p *Processor) SaveProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMessageType: //登录
		var up process.UserProcess
		up = process.UserProcess{
			Conn: p.Conn,
		}
		up.SaveProcessLogin(mes)
		return nil
	case message.RegisterMessageType: //注册
		var up process.UserProcess
		up = process.UserProcess{
			Conn: p.Conn,
		}
		up.SaveProcessRegister(mes)
		return nil
	case message.UserListMessageType: //用户列表
		var up process.UserProcess
		up = process.UserProcess{
			Conn: p.Conn,
		}
		up.SaveProcessUserList(mes)
		return nil
	case message.SmsMessageType: //发送消息
		var up process.SmsProcess
		up = process.SmsProcess{}
		up.Message(mes)
		return nil
	default:
		fmt.Println("消息类型不正确")
	}
	return nil
}

func (p *Processor) ProcessConn() (err error) {
	for {
		mes, err := utils.ReadMessage(p.Conn)
		if err != nil {
			return err
		}
		p.SaveProcessMes(mes)
	}
	return err
}
