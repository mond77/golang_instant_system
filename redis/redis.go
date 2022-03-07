package redis

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     8,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")

		},
	}

}
func VerifyUser(name string, psw string) bool {
	conn := pool.Get()
	psw01, err := redis.String(conn.Do("get", name))
	if err != nil {
		fmt.Println("conn get err:", err)
		return false
	}
	if psw != psw01 {
		fmt.Println("user name is wrong,please reset!")
		return false
	}
	fmt.Println("登陆成功")

	return true

}
