package mysql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zzlpeter/aps-go/libs/tomlc"
	"log"
	"sync"
)

var dbMap = map[string]*gorm.DB{}
var once sync.Once

func getDbConn() {
	once.Do(func() {
		makeDbConn()
	})
}

func makeDbConn() {
	dbConf := tomlc.Config{}.MysqlConfS()
	for alias, conf := range dbConf {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds",
			conf["username"], conf["password"], conf["host"],
			conf["port"], conf["db"], conf["timeout"])
		db, err := gorm.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("connect mysql: %v err: %v", alias, err.Error())
		}
		maxConn := conf["max_con"].(int)
		db.DB().SetMaxOpenConns(maxConn)
		maxIdle := conf["max_idle"].(int)
		db.DB().SetMaxIdleConns(maxIdle)
		dbMap[alias] = db
	}
}

func GetDbConn(db string) (*gorm.DB, error) {
	getDbConn()
	if val, ok := dbMap[db]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("db alias: %s is missed, please ensure you make it in conf.toml", db))
}

func init() {
	getDbConn()
}
/*
type Task struct {
	ID 		int
	Status  string
}

func (t Task) TableName() string {
	return "task"
}

func testQuery() {
	db, _ := mysql.GetDbConn("default")
	task := Task{}
	db.Select([]string{"id", "status"}).First(&task)
	fmt.Println(task.ID, task.Status)
}
*/