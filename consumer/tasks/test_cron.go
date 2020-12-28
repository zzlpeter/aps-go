package tasks

import (
	"fmt"
	"github.com/zzlpeter/aps-go/models"
	"time"
)

func TestCron(p models.TaskQueueStruct) {
	time.Sleep(time.Second * 2)
	fmt.Println(time.Now(), "TestCron", p)
	panic("test msg")
}
