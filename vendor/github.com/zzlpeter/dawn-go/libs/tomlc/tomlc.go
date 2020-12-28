package tomlc

import (
	"github.com/BurntSushi/toml"
	"github.com/zzlpeter/dawn-go/libs/utils"
	"log"
	"sync"
)


type tomlConfig struct {
	Mysql	map[string]mysql      `toml:"mysql"`
	Redis 	map[string]redis      `toml:"redis"`
	Basic   map[string]interface{} `toml:"basic"`
}

var (
	cfg  tomlConfig
	once sync.Once
)

type mysql struct {
	Host 		string 		`json:"host" map:"host" toml:"host"`
	Port 		int			`json:"port" map:"port" toml:"port"`
	Username  	string		`json:"username" map:"username" toml:"username"`
	Password 	string		`json:"password" map:"password" toml:"password"`
	Db 			string		`json:"db" map:"db" toml:"db"`
	MaxCon		int 		`json:"max_con" map:"max_con" toml:"max_con"`
	MaxIdle 	int 		`json:"max_idle" map:"max_idle" toml:"max_idle"`
	Timeout 	int 		`json:"timeout" map:"timeout" toml:"timeout"`
}

type redis struct {
	Host 		string		`json:"host" map:"host" toml:"host"`
	Port		int			`json:"port" map:"port" toml:"port"`
	Db 			int			`json:"db" map:"db" toml:"db"`
	Password 	string		`json:"password" map:"password" toml:"password"`
	MaxIdle 	int 		`json:"max_idle" map:"max_idle" toml:"max_idle"`
	MaxActive 	int 		`json:"max_active" map:"max_active" toml:"max_active"`
	IdleTimeout	int 		`json:"idle_timeout" map:"idle_timeout" toml:"idle_timeout"`
}

func getFilePath() string {
	return "/Users/zhangzhiliang/go/src/github.com/zzlpeter/dawn-go/conf/conf.toml"
}

func readConfig() {
	fileConf := getFilePath()
	_, err := toml.DecodeFile(fileConf, &cfg)
	if err != nil {
		log.Fatalf("read conf.toml fails with error: %v", err.Error())
	}
	log.Println("read config init...")
}

func getConfig() *tomlConfig {
	once.Do(func() {
		readConfig()
	})
	return &cfg
}

type Config struct {
}

func (c Config) MysqlConf(alias string) map[string]interface{} {
	conf := getConfig()
	cS := conf.Mysql[alias]
	m := utils.Struct2Map(cS)
	return m
}

func (c Config) MysqlConfS() map[string]map[string]interface{} {
	conf := getConfig()
	mysqlMap := make(map[string]map[string]interface{})
	for k, v := range conf.Mysql {
		m := utils.Struct2Map(v)
		mysqlMap[k] = m
	}
	return mysqlMap
}

func (c Config) RedisConf(alias string) map[string]interface{} {
	conf := getConfig()
	cS := conf.Redis[alias]
	m := utils.Struct2Map(cS)
	return m
}

func (c Config) RedisConfS() map[string]map[string]interface{} {
	conf := getConfig()
	redisMap := make(map[string]map[string]interface{})
	for k, v := range conf.Redis {
		m := utils.Struct2Map(v)
		redisMap[k] = m
	}
	return redisMap
}

func (c Config) BasicConf() map[string]interface{} {
	conf := getConfig()
	return conf.Basic
}

func init() {
	getConfig()
}