package main

import (
	"flag"
	"fmt"
	"github.com/zzlpeter/aps-go/consumer"
	aps_log "github.com/zzlpeter/aps-go/libs/log"
	_ "github.com/zzlpeter/aps-go/libs/log"
	_ "github.com/zzlpeter/aps-go/libs/mysql"
	_ "github.com/zzlpeter/aps-go/libs/redis"
	"github.com/zzlpeter/aps-go/libs/utils"
	_ "github.com/zzlpeter/aps-go/libs/utils"
	"github.com/zzlpeter/aps-go/producer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 生产者
func apsProducer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	prd := producer.NewProducerManager()
	prd.ProduceJobs2Queue()
	prd.Start()
	for range ticker.C {
		prd.ProduceJobs2Queue()
	}
}

// 消费者
func apsConsumer() {
	con := consumer.NewConsumerManager()
	con.ConsumerMsgFromQueue()
}

var action = flag.String("action", "consumer", "请输入操作类型: consumer or producer")

// 信号处理逻辑
func signalTerm(c chan os.Signal) {
	<- c
	utils.Environ{}.Set("IS_KILLED", "IS_KILLED")
	aps_log.LogRecord(*action, "", aps_log.WARNING, fmt.Sprintf("%v收到term信号，程序即将退出", *action))
	if *action == "consumer" {
		time.Sleep(time.Second * utils.ConsumerTimeout)
	} else {
		time.Sleep(time.Second * utils.ProducerTimeout)
	}

	os.Exit(1)
}

func main() {
	// 捕获TERM信号 - 优雅重启
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	go signalTerm(ch)

	// 解析参数、启动程序
	flag.Parse()
	if *action == "consumer" {
		apsConsumer()
	} else if *action == "producer" {
		apsProducer()
	} else {
		log.Fatalf("action: %v 参数无效，可选值为consumer/producer", *action)
	}
}