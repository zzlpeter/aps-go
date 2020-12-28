package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/zzlpeter/aps-go/libs/tomlc"
	"log"
	"sync"
	"time"
)

type pool struct {
	redis.Pool
}

// get val
func (p *pool) GetVal(key string) (interface{}, error) {
	rc := p.Get()
	defer rc.Close()
	rly, err := rc.Do("GET", key)
	return rly, err
}

// set val
func (p *pool) SetVal(key string, args ...interface{}) (interface{}, error) {
	rc := p.Get()
	defer rc.Close()
	rly, err := rc.Do("SET", key, args)
	return rly, err
}

// del val
func (p *pool) DeleteVal(key string) (interface{}, error) {
	rc := p.Get()
	defer rc.Close()
	rly, err := rc.Do("DEL", key)
	return rly, err
}

var redisPollMap = map[string]*pool{}
var once sync.Once

func getPool() {
	once.Do(func() {
		makeRedisPool()
	})
}

func makeRedisPool() {
	redisConf := tomlc.Config{}.RedisConfS()
	for alias, conf := range redisConf {
		var client *pool
		client = &pool{
			redis.Pool{
				MaxIdle:     conf["max_idle"].(int),
				MaxActive:   conf["max_active"].(int),
				IdleTimeout: time.Duration(conf["idle_timeout"].(int)),
				Wait:        true,
				Dial:        dial(alias, conf),
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			},
		}

		// 校验连接是否正确建立
		ch := make(chan bool)
		go testConnectRedis(ch, client)
		select {
		case <- time.After(time.Second * 5):
			log.Fatalf("connect redis: dial tcp %v:%v connect timeout", conf["host"], conf["port"])
		case <- ch:
		}

		redisPollMap[alias] = client
	}
}

func testConnectRedis(ch chan bool, rds *pool) {
	r := rds.Get()
	defer r.Close()
	ch <- true
}

func dial(alias string, conf map[string]interface{}) func() (redis.Conn, error) {
	f := func() (redis.Conn, error) {
		con, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", conf["host"], conf["port"]))
		if err != nil {
			log.Fatalf("connect redis: %v err: %v", alias, err.Error())
		}
		if pwd := conf["password"].(string); pwd != "" {
			if _, err := con.Do("AUTH", pwd); err != nil {
				log.Fatalf("auth redis: %v err: %v", alias, err.Error())
			}
		}
		if _, err := con.Do("SELECT", conf["db"]); err != nil {
			log.Fatalf("db select redis: %v err: %v", alias, err.Error())
		}
		return con, nil
	}
	return f
}

func GetRedisPool(alias string) *pool {
	getPool()
	return redisPollMap[alias]
}

func init() {
	getPool()
}

/*
func test() {
	rc := GetRedisPool("default").Get()
	defer rc.Close()
	rst, err := rc.Do("get", "name")
	log.Println(rst, err)
}
*/