package main

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type lenResult struct {
	Key string
	Len int64
}

func CheckLen(dsn RedisDSN) (map[string]int64, error) {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:        dsn.Addr,
		Password:    dsn.Password, // no password set
		DB:          dsn.DB,       // use default DB
		ReadTimeout: 10 * time.Second,
		PoolSize:    10,
		DialTimeout: 10 * time.Second,
	})

	defer redisClient.Close()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)

	// 判断是否key包含*
	if strings.Contains(dsn.Key, "*") {
		// 使用Scan函数进行迭代遍历
		var cursor uint64
		var keys []string
		ch := make(chan lenResult, 1000)
		stopCh := make(chan struct{})
		waitGroup := sync.WaitGroup{}

		go func() {
			for v := range ch {
				result[v.Key] = v.Len
			}
			stopCh <- struct{}{}
		}()

		// 走 scan match 方式
		for {
			var err error
			keys, cursor, err = redisClient.Scan(ctx, cursor, dsn.Key, 1000).Result()
			if err != nil {
				continue
			}

			// 处理匹配的键名
			for _, key := range keys {
				waitGroup.Add(1)
				go func(key string) {
					vLen, err := llen(context.Background(), redisClient, key)

					if err != nil {
						waitGroup.Done()
						return
					}

					ch <- lenResult{
						Key: key,
						Len: vLen,
					}

					waitGroup.Done()
				}(key)
			}

			// 如果游标为0，表示迭代结束
			if cursor == 0 {
				waitGroup.Wait()
				close(ch)
				break
			}
		}

		<-stopCh
	} else {
		vLen, err := llen(ctx, redisClient, dsn.Key)
		if err != nil {
			return nil, err
		}

		result[dsn.Key] = vLen
	}

	return result, nil
}

func llen(ctx context.Context, redisClient *redis.Client, key string) (int64, error) {
	// 判断key类型
	keyType, err := redisClient.Type(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	var vLen int64

	if keyType == "list" {
		vLen, err = redisClient.LLen(ctx, key).Result()
		if err != nil {
			return 0, err
		}
	} else if keyType == "zset" {
		vLen, err = redisClient.ZCard(ctx, key).Result()
		if err != nil {
			return 0, err
		}
	}

	return vLen, err
}
