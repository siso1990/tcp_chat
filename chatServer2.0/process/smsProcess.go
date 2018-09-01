package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"go_code/chat2.0/common/message"
)

type SmsProcess struct {
}

func (s *SmsProcess) Message(mes *message.Message) (err error) {
	var sms message.SmsMessage
	err = json.Unmarshal([]byte(mes.Data), &sms)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &sms)", err)
		return
	}
	data, err := json.Marshal(*mes)
	if err != nil {
		fmt.Println("json.Marshal(mes)", err)
		return
	}

	for i, v := range OnlineUsers {
		//发送文件大小
		var buf [8096]byte
		if sms.Addressee == 0 || sms.Addressee == i {
			pkgLen := uint32(len(data))
			binary.BigEndian.PutUint32(buf[:4], pkgLen)
			n, err := v.Conn.Write(buf[:4])
			if err != nil || n != 4 {
				fmt.Println("conn.Write(buf)", err)
				return err
			}
			//发送mes
			_, err = v.Conn.Write(data)
			if err != nil {
				fmt.Println("conn.Write(data)", err)
				return err
			}
		}
	}
	return
}
