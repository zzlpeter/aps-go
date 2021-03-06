package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	aps_log "github.com/zzlpeter/aps-go/libs/log"
	"github.com/zzlpeter/aps-go/libs/mysql"
	aps_redis "github.com/zzlpeter/aps-go/libs/redis"
	"github.com/zzlpeter/aps-go/libs/tomlc"
	"github.com/zzlpeter/aps-go/libs/utils"
	"github.com/zzlpeter/aps-go/models"
	"runtime/debug"
	"time"
)


var queue = tomlc.Config{}.BasicConf()["task_redis_queue"].(string)
type consumerStruct struct {}

// callbackFunc 查找回调函数
func (c *consumerStruct) callbackFunc(t models.TaskQueueStruct) func(t models.TaskQueueStruct) {
	f, ok := TaskKeyFuncMapper[t.ExecuteFunc]
	if ok {
		return f
	}
	return nil
}

// do 执行业务逻辑
func (c *consumerStruct) do(t models.TaskQueueStruct) {
	// 查找回调方法
	callback := c.callbackFunc(t)
	if callback == nil {
		aps_log.LogRecord("consumer", t.TraceId, aps_log.ERROR, fmt.Sprintf("未发现%v对应的注册函数", t.ExecuteFunc))
		return
	}

	// 更新子任务状态
	db, _ := mysql.GetDbConn("default")
	db.Model(&models.TaskExecute{}).Where("id = ?", t.SubTaskId).Update("status", "doing")

	// 捕获函数可能出现的异常
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

	// 执行函数
	aps_log.LogRecord("consumer", t.TraceId, aps_log.INFO, fmt.Sprintf("开始执行函数:%v", t.ExecuteFunc))
	callback(t)
	aps_log.LogRecord("consumer", t.TraceId, aps_log.INFO, fmt.Sprintf("执行函数结束:%v", t.ExecuteFunc))
}

// ConsumerMsgFromQueue 从队列消费消息
func (c *consumerStruct) ConsumerMsgFromQueue() {
	conn := aps_redis.GetRedisPool("default").Get()
	defer conn.Close()
	aps_log.LogRecord("consumer", "", aps_log.INFO, "消费worker开始工作")

	for {
		// 判断环境变量是否为 IS_KILLED
		isKilled := utils.Environ{}.Get("IS_KILLED")
		if isKilled == "IS_KILLED" {
			aps_log.LogRecord("consumer", "", aps_log.WARNING, "consumer环境变量为IS_KILLED,停止消费消息")
			// 此处不应该主动退出、由main函数统一处理
			time.Sleep(time.Second * utils.ConsumerTimeout)
		}

		// 阻塞读取消息
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

		// 开启新的协程执行业务逻辑
		go func(_t models.TaskQueueStruct) {
			c.do(_t)
		}(t)
	}
}

func NewConsumerManager() *consumerStruct {
	return &consumerStruct{}
}