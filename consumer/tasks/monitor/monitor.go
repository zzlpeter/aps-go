package monitor

import (
	"fmt"
	"github.com/zzlpeter/aps-go/libs/mysql"
	"github.com/zzlpeter/aps-go/libs/utils"
	"github.com/zzlpeter/aps-go/models"
)

func TasksMonitor(t models.TaskQueueStruct) {
	// 获取所有有效的任务
	tasks := []models.Task{}
	db, _ := mysql.GetDbConn("default")
	db.Select([]string{"id", "task_key", "extra"}).Where("is_valid = 1").Find(&tasks)
	// 循环判断每个任务执行情况
	alarm := []string{"任务key - 延迟s"}
	for _, task := range tasks {
		sub := models.TaskExecute{}
		db.Select([]string{"id", "create_at"}).Where("task_id = ?", task.Id).Where("status in (?)", []string{"doing", "fail", "success"}).Order("id desc").Limit(1).Find(&sub)
		if sub.Id == 0 {
			continue
		}
		// 获取延迟时间
		delaySeconds := 3600
		if t, ok := task.Extra["delay"]; ok {
			if v, o := t.(int); o {
				delaySeconds = v
			}
		}
		diffTime := utils.StampSecond() - sub.CreateAt.Unix()
		if diffTime > int64(delaySeconds) {
			alarm = append(alarm, fmt.Sprintf("%v - %v", task.TaskKey, diffTime))
		}
	}
	if len(alarm) > 1 {
		// alarm
	}
}