package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	 log "github.com/sirupsen/logrus"
	_ "fmt"
)

var rdb *redis.ClusterClient
var ctx = context.Background()

func init() {
	node1 := os.Getenv("redis_node_1")
	node2 := os.Getenv("redis_node_2")
	node3 := os.Getenv("redis_node_3")
	node4 := os.Getenv("redis_node_4")
	node5 := os.Getenv("redis_node_5")
	node6 := os.Getenv("redis_node_6")
	password := os.Getenv("redis_pass")

	addrs := []string{node1, node2, node3, node4, node5, node6}

	conn := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})

	err := conn.Ping(ctx).Err()
	if err != nil {
		log.Fatal("[Redis] Failed to connect to cluster: ",err)
	} else {
		rdb = conn
		log.Info("[Redis] Cluster connected successfully (6 nodes)")
	}
}

func GetRedis() *redis.ClusterClient {
	return rdb
}

