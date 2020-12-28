package models

import (
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"github.com/zzlpeter/aps-go/libs/mysql"
	"github.com/zzlpeter/aps-go/libs/utils"
	"time"
)

type DateTimeStruct struct {
	CreateAt 	time.Time 	`json:"create_at" comment:"创建时间" gorm:"column:create_at;default:null"`
	UpdateAt 	time.Time	`json:"update_at" comment:"更新时间" gorm:"column:update_at;default:null"`
}

type Task struct {
	DateTimeStruct
	Id 			uint 		`json:"id" comment:"主键自增" gorm:"primary_key"`
	TaskKey 	string 		`json:"task_key" comment:"任务唯一key" gorm:"index"`
	Desc 		string 		`json:"desc" comment:"任务描述"`
	ExecuteFunc string		`json:"execute_func" comment:"执行方法"`
	Spec 		string		`json:"spec" comment:"执行时间"`
	Params 		GORMJsonMapper		`json:"params" comment:"执行参数"`
	IsValid 	bool		`json:"is_valid" comment:"是否有效"`
	Status 		string		`json:"status" comment:"任务执行状态ready/doing"`
	Extra 		GORMJsonMapper		`json:"extra" comment:"额外信息(json格式)"`
}

func (t *Task) TableName() string {
	return "task"
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	// 校验extra是否有值
	if t.Extra == nil {
		t.Extra = map[string]interface{}{}
	}
	return nil
	// 校验trigger是否为cron/date/interval
	//switch t.Trigger {
	//case "cron":
	//	return t.isCronValid()
	//case "interval":
	//	return t.isIntervalValid()
	//case "date":
	//	return t.isDateValid()
	//default:
	//	err := fmt.Sprintf("Trigger: %v is invalid", t.Trigger)
	//	return errors.New(err)
	//}
}

func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}

func (t Task) isCronValid() error {
	_, err := cron.ParseStandard(t.Spec)
	return err
}

func (t Task) isIntervalValid() error {
	_, err := utils.String2Int(t.Spec)
	return err
}

func (t Task) isDateValid() error {
	_, err := time.Parse("2006-01-02 15:04:05", t.Spec)
	return err
}

type TaskExecute struct {
	DateTimeStruct
	Id 			uint 		`json:"id" comment:"主键自增" gorm:"primary_key"`
	TaskId 		uint		`json:"task_id" comment:"任务ID" gorm:"index"`
	TraceId 	string 		`json:"trace_id" comment:"traceID"`
	Status 		string 		`json:"status" comment:"执行状态,初始值doing"`
	Extra 		GORMJsonMapper 		`json:"extra" comment:"额外信息"`
}

func (t *TaskExecute) TableName() string {
	return "task_execute"
}

func (t *TaskExecute) UpdateExt(extra map[string]interface{}, id uint) error {
	var ext = make(map[string]interface{})
	db, _ := mysql.GetDbConn("default")
	if t.Id != 0 {
		for k := range t.Extra {
			ext[k] = t.Extra[k]
		}
	} else {
		sub := TaskExecute{}
		db.Model(&TaskExecute{}).Where("id = ?", id).Limit(1).Find(&sub)
		for k := range sub.Extra {
			ext[k] = sub.Extra[k]
		}
	}
	for k := range extra {
		ext[k] = extra[k]
	}
	if err := db.Model(&TaskExecute{}).Where("id = ?", id).Update("extra", ext).Error; err != nil {
		return err
	}
	return nil
}

type TaskQueueStruct struct {
	TraceId 	string 					`json:"trace_id"`
	SubTaskId 	uint					`json:"sub_task_id"`
	Params 		map[string]interface{}	`json:"params"`
	ExecuteFunc string					`json:"execute_func"`
}