package main

import (
	"fmt"
	"net"
	"time"

	"go_code/chat2.0/chatServer2.0/model"
	"go_code/chat2.0/chatServer2.0/process"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

//初始化数据库连接
func init() {
	InitPool(16, 0, 300*time.Second, "localhost:6379")
	model.MyUserDao = model.NewUserDao(pool)
}

//管理连接
func ConnProcess(conn net.Conn) {
	defer func() {
		conn.Close()
		err := recover()
		if err != nil {
			fmt.Println("用户%v已经断开连接", conn.RemoteAddr())
		}
		for i, v := range process.OnlineUsers {
			if v.Conn == conn {
				delete(process.OnlineUsers, i)
				continue
			}
		}
	}()
	p := &Processor{conn}
	p.ProcessConn()
}

//建立连接池
func InitPool(maxIdle, maxActive int, IdleTimeout time.Duration, address string) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: IdleTimeout,
		Dial: func() (redis.Conn, error) { // 初始化链接的代码， 链接哪个ip的redis
			return redis.Dial("tcp", address)
		},
	}
}

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("开启监听失败, 试图开启 8889 端口", err)
		panic(err)
	}
	for {
		fmt.Println("开始监听端口")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("发生一次失败连接", err)
			panic(err)
		}
		fmt.Printf("%v成功连接服务器", conn.RemoteAddr())
		go ConnProcess(conn)
	}
}
