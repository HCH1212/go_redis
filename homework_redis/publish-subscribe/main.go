// redis实现一个发布订阅模型
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 创建一个Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	// 订阅一个频道
	pubsub := rdb.Subscribe(ctx, "mychannel")
	defer pubsub.Close()

	// 创建一个接收消息的通道
	ch := pubsub.Channel()

	// 启动一个goroutine来处理接收到的消息
	//在这个goroutine中，通过循环从通道ch中读取消息，并打印出消息的频道和内容
	go func() {
		for msg := range ch {
			fmt.Println("Received message:", msg.Channel, msg.Payload)
		}
	}()

	// 发布一些消息到频道
	for i := 0; i < 10; i++ {
		err := rdb.Publish(ctx, "mychannel", fmt.Sprintf("message %d", i)).Err()
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	// 等待一段时间以便接收所有消息
	time.Sleep(5 * time.Second)
}

//https://www.tizi365.com/archives/306.html
