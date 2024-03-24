package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// 添加用户分数
func addScore(rdb *redis.Client, user string, score int) {
	rdb.ZAdd(ctx, "leaderboard", &redis.Z{Score: float64(score), Member: user})
}

// 获取排行榜前n名
func getTopUsers(rdb *redis.Client, limit int) []string {
	result, err := rdb.ZRevRangeByScoreWithScores(ctx, "leaderboard", &redis.ZRangeBy{
		Min:    "0",
		Max:    "+inf",
		Offset: 0,
		Count:  int64(limit),
	}).Result()
	if err != nil {
		fmt.Println("rdb.ZRevRangeByScoreWithScores err=", err)
	}

	users := make([]string, len(result))
	for i, z := range result {
		users[i] = fmt.Sprintf("%s: %f", z.Member, z.Score)
	}
	return users
}

func main() {
	// 连接到Redis服务器
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 如果没有密码，留空
		DB:       0,  // 默认数据库
	})

	// 添加用户分数
	addScore(rdb, "user1", 100)
	addScore(rdb, "user2", 200)
	addScore(rdb, "user3", 300)
	addScore(rdb, "user4", 900)
	addScore(rdb, "user5", 200)
	addScore(rdb, "user6", 9900)
	addScore(rdb, "user7", 500)
	addScore(rdb, "user8", 1400)
	addScore(rdb, "user9", 600)
	addScore(rdb, "user10", 300)

	// 获取排行榜前n名
	n := 6
	topUsers := getTopUsers(rdb, 6)
	fmt.Printf("排行榜前%d名:", n)
	fmt.Println(topUsers)
}

//https://www.tizi365.com/archives/304.html
