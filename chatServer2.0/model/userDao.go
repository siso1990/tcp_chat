package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"go_code/chat2.0/common/message"
)

var MyUserDao *UserDao

const (
	userKey string = "users"
)

type UserDao struct {
	Pool *redis.Pool
}

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return
}

func (ud *UserDao) GetUserById(conn redis.Conn, id int) (user *message.User, err error) {
	result, err := redis.String(conn.Do("HGet", userKey, id))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &message.User{}
	err = json.Unmarshal([]byte(result), user)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (ud *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {
	conn := ud.Pool.Get()
	defer conn.Close()
	user, err = ud.GetUserById(conn, userId)
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
		fmt.Println("ud.GetUserById(conn,userId)", err)
		return
	}
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
	}
	return
}

func (ud *UserDao) Register(userId int, userPwd string,
	userName string) (err error) {
	conn := ud.Pool.Get()
	defer conn.Close()

	user, err := ud.GetUserById(conn, userId)
	if err != nil {
		if err == ERROR_USER_NOTEXISTS {
			var user message.User
			user.UserId = userId
			user.UserPwd = userPwd
			user.UserName = userName
			data, err := json.Marshal(user)
			if err != nil {
				fmt.Println("json.Marshal(user) err=", err)
				return err
			}
			conn.Do("Hset", userKey, user.UserId, data)
			return nil
		} else {
			return
		}
	} else if user.UserId == userId {
		err = ERROR_USER_EXISTS
	}
	return
}
