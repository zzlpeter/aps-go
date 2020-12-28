package producer

import (
	"fmt"
	"github.com/robfig/cron/v3"
	aps_log "github.com/zzlpeter/aps-go/libs/log"
	"github.com/zzlpeter/aps-go/libs/mysql"
	"github.com/zzlpeter/aps-go/libs/redis"
	"github.com/zzlpeter/aps-go/libs/tomlc"
	"github.com/zzlpeter/aps-go/libs/utils"
	"github.com/zzlpeter/aps-go/models"
	"sync"
	"time"
)

type _job struct {
	Gid 	cron.EntryID
	Spec 	string
}

// 自定义 withLogger
type apsCronZap struct {}

func (a apsCronZap) Info(msg string, kvs ...interface{}) {
	aps_log.LogRecord("producer", "", aps_log.INFO, msg, kvs...)
}
func (a apsCronZap) Error(err error, msg string, kvs ...interface{}) {
	kvs = append(kvs, "error", err)
	aps_log.LogRecord("producer", "", aps_log.ERROR, msg, kvs...)
}

var cronManager = cron.New(cron.WithLogger(apsCronZap{}))
var queue = tomlc.Config{}.BasicConf()["task_redis_queue"].(string)
type producerStruct struct {
	jobsSpecMapper 	map[uint]_job
	mu 				sync.Mutex
}

// 发布消息到队列
func (p *producerStruct) msg2Q(t models.Task) {
	fmt.Println("msg2q", time.Now())
	aps_log.LogRecord("producer", "", aps_log.INFO, fmt.Sprintf("方法:%v 开始生成消息", t.ExecuteFunc))
	rds := redis.GetRedisPool("default").Get()
	defer rds.Close()

	tid := utils.GenUUID()
	err, subTaskId := p.addSubTask(t, tid)
	if err != nil {
		aps_log.LogRecord("producer", tid, aps_log.ERROR, "插入子任务记录失败", "err", err.Error(), "task_key", t.TaskKey, "spec", t.Spec)
		return
	}

	msgStruct := models.TaskQueueStruct{
		TraceId: tid,
		SubTaskId: subTaskId,
		Params: t.Params,
		ExecuteFunc: t.ExecuteFunc,
	}
	js, err := utils.Struct2Json(msgStruct)
	if err != nil {
		aps_log.LogRecord("producer", tid, aps_log.ERROR, "序列化结构体失败", "err", err.Error(), "task_key", t.TaskKey, "spec", t.Spec)
		return
	}
	_, err = rds.Do("LPUSH", queue, js)
	if err != nil {
		aps_log.LogRecord("producer", tid, aps_log.ERROR, "producer发布消息失败", "err", err.Error(), "task_key", t.TaskKey, "spec", t.Spec)
		// alarm
	}
}

// 插入执行记录到子任务表
func (p *producerStruct) addSubTask(t models.Task, tid string) (error, uint) {
	insert := models.TaskExecute{TaskId: t.Id, TraceId: tid, Extra: map[string]interface{}{}, Status: "todo"}
	db, _ := mysql.GetDbConn("default")
	if err := db.Create(&insert).Error; err != nil {
		return err, 0
	}
	return nil, insert.Id
}

// 添加任务
func (p *producerStruct) addJob(t models.Task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	spec := t.Spec
	jid, err := cronManager.AddFunc(spec, func() {
		p.msg2Q(t)
	})
	if err != nil {
		aps_log.LogRecord("producer", "", aps_log.ERROR, "添加任务失败", "err", err.Error(), "task_key", t.TaskKey, "spec", t.Spec)
		// alarm
	}
	p.jobsSpecMapper[t.Id] = _job{jid, t.Spec}
}

func (p *producerStruct) ProduceJobs2Queue() {
	db, _ := mysql.GetDbConn("default")
	tasks := []models.Task{}
	db.Select([]string{"id", "spec", "params", "is_valid", "execute_func"}).Find(&tasks)
	for _, tk := range tasks {
		job, jok := p.jobsSpecMapper[tk.Id]
		// 任务已无效 -> 「任务移除」
		if !tk.IsValid {
			if jok {
				aps_log.LogRecord("producer", "", aps_log.INFO, "移除任务", "task_key", tk.TaskKey)
				cronManager.Remove(job.Gid)
			}
			continue
		}
		// 任务已存在 && 调度频次已改变 ->「先移除再新增」
		if jok {
			if job.Spec != tk.Spec {
				aps_log.LogRecord("producer", "", aps_log.INFO, "更新任务调度频次", "task_key", tk.TaskKey)
				cronManager.Remove(job.Gid)
				p.addJob(tk)
			}
		// 任务不存在 ->「直接新增」
		} else {
			aps_log.LogRecord("producer", "", aps_log.INFO, "新增任务", "task_key", tk.TaskKey)
			p.addJob(tk)
		}
	}
}

func (p *producerStruct) Start() {
	cronManager.Start()
}

func (p *producerStruct) Stop() {
	cronManager.Stop()
}

func NewProducerManager() *producerStruct {
	return &producerStruct{jobsSpecMapper: make(map[uint]_job)}
}