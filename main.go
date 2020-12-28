package main

import (
	"flag"
	"github.com/zzlpeter/aps-go/consumer"
	_ "github.com/zzlpeter/aps-go/libs/log"
	_ "github.com/zzlpeter/aps-go/libs/mysql"
	_ "github.com/zzlpeter/aps-go/libs/redis"
	_ "github.com/zzlpeter/aps-go/libs/utils"
	"github.com/zzlpeter/aps-go/producer"
	"log"
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

func main() {
	flag.Parse()
	if *action == "consumer" {
		apsConsumer()
	} else if *action == "producer" {
		apsProducer()
	} else {
		log.Fatalf("action: %v 参数无效，可选值为consumer/producer", *action)
	}
}