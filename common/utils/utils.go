package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"go_code/chat2.0/common/message"
)

//读取方法
func ReadMessage(conn net.Conn) (mes *message.Message, err error) {
	var buf [8096]byte
	_, err = conn.Read(buf[:4])
	if err != nil {
		fmt.Println("读取信息时出错", err)
		return
	}
	//读取mes

	pkgLen := binary.BigEndian.Uint32(buf[:4])
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Printf("读取消息时发生错误", err)
		return
	}
	err = json.Unmarshal(buf[:pkgLen], &mes)
	return
}

//写入方法
func WritMessage(conn net.Conn, mes *message.Message) (err error) {
	//序列化mes
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("解析消息时发生错误", err)
		return
	}
	//发送文件大小
	var buf [8096]byte
	pkgLen := uint32(len(data))
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	n, err := conn.Write(buf[:4])
	if err != nil || n != 4 {
		fmt.Println("conn.Write(buf)", err)
		return
	}
	//发送mes
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("读取数据时发生错误", err)
		return
	}
	return
}
