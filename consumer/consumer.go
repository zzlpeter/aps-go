package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/zzlpeter/aps-go/consumer/tasks"
	aps_log "github.com/zzlpeter/aps-go/libs/log"
	"github.com/zzlpeter/aps-go/libs/mysql"
	aps_redis "github.com/zzlpeter/aps-go/libs/redis"
	"github.com/zzlpeter/aps-go/libs/tomlc"
	"github.com/zzlpeter/aps-go/models"
	"runtime/debug"
)

// 新增任务需要在此进行注册
var taskKeyFuncMapper = map[string]func(s models.TaskQueueStruct){
	"TestCron": tasks.TestCron,
}

// please do not change below !!!
var queue = tomlc.Config{}.BasicConf()["task_redis_queue"].(string)
type consumerStruct struct {}

func (c *consumerStruct) callbackFunc(t models.TaskQueueStruct) func(t models.TaskQueueStruct) {
	f, ok := taskKeyFuncMapper[t.ExecuteFunc]
	if ok {
		return f
	}
	return nil
}

func (c *consumerStruct) do(t models.TaskQueueStruct) {
	callback := c.callbackFunc(t)
	if callback == nil {
		aps_log.LogRecord("consumer", t.TraceId, aps_log.ERROR, fmt.Sprintf("未发现%v对应的注册函数", t.ExecuteFunc))
		return
	}
	db, _ := mysql.GetDbConn("default")
	db.Model(&models.TaskExecute{}).Where("id = ?", t.SubTaskId).Update("status", "doing")
	defer func() {
		var status = "success"
		ext := map[string]interface{}{}
		r := recover()
		if r != nil {
			status = "fail"
			ext["error"] = string(debug.Stack())
			aps_log.LogRecord("consumer", t.TraceId, aps_log.ERROR, fmt.Sprintf("执行函数:%v 异常", t.ExecuteFunc), "err", string(debug.Stack()))
		}
		db.Model(&models.TaskExecute{}).Where("id = ?", t.SubTaskId).Update("status", status)
		if len(ext) != 0 {
			(&models.TaskExecute{}).UpdateExt(ext, t.SubTaskId)
		}
	}()
	aps_log.LogRecord("consumer", t.TraceId, aps_log.INFO, fmt.Sprintf("开始执行函数:%v", t.ExecuteFunc))
	callback(t)
	aps_log.LogRecord("consumer", t.TraceId, aps_log.INFO, fmt.Sprintf("执行函数结束:%v", t.ExecuteFunc))
}

func (c *consumerStruct) ConsumerMsgFromQueue() {
	conn := aps_redis.GetRedisPool("default").Get()
	defer conn.Close()
	aps_log.LogRecord("consumer", "", aps_log.INFO, "消费worker开始工作")

	for {
		reply, err := redis.Values(conn.Do("BRPOP", queue, 0))
		if err != nil {
			aps_log.LogRecord("consumer", "", aps_log.ERROR, "消费worker获取Redis消息异常", "err", err.Error())
			// alarm
			continue
		}
		var t = models.TaskQueueStruct{}
		if err := json.Unmarshal(reply[1].([]byte), &t); err != nil {
			aps_log.LogRecord("consumer", "", aps_log.ERROR, "反序列化json异常", "err", err.Error())
			// alarm
			continue
		}
		go func(_t models.TaskQueueStruct) {
			c.do(_t)
		}(t)
	}
}

func NewConsumerManager() *consumerStruct {
	return &consumerStruct{}
}