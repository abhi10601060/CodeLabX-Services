package rds

import "github.com/redis/go-redis/v9"

const (
	url = "redis://Abhi1060:Abhi1060@localhost:6379/0?protocol=3"
)

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})
}