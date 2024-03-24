package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	lockKey := "my_lock"
	lockValue := "lock_value"
	lockTimeout := 10 * time.Second

	// 尝试获取锁
	success, err := rdb.SetNX(ctx, lockKey, lockValue, lockTimeout).Result()
	if err != nil {
		panic(err)
	}

	if success {
		fmt.Println("获取锁成功")
		// 执行业务逻辑
		// ...

		// 释放锁
		unlockScript := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
		`)
		_, err := unlockScript.Run(ctx, rdb, []string{lockKey}, lockValue).Result()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("获取锁失败")
	}
}

//https://blog.csdn.net/cljdsc/article/details/123385538?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522171128652716800226510246%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=171128652716800226510246&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~sobaiduend~default-1-123385538-null-null.142^v99^pc_search_result_base9&utm_term=golang%20redis%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81&spm=1018.2226.3001.4187
