package consumer

import (
	"github.com/zzlpeter/aps-go/consumer/tasks"
	"github.com/zzlpeter/aps-go/consumer/tasks/monitor"
	"github.com/zzlpeter/aps-go/models"
)

// 新增任务需要在此进行添加
var TaskKeyFuncMapper = map[string]func(s models.TaskQueueStruct){
	"TestCron": tasks.TestCron,
	"TasksMonitor": monitor.TasksMonitor,
}
