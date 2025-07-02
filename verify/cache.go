package verify

import (
	"context"
	"encoding/json"
	"fmt"
	logger "checker-middleware/pkg/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

// 获取 redis 客户端
func getRedisClient(cfg CacheConfig) (redis.UniversalClient, error) {
	opts := &redis.UniversalOptions{
		Addrs:    []string{cfg.Host + ":" + fmt.Sprintf("%d", cfg.Port)},
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	switch cfg.Mode {
	case "sentinel":
		opts.Addrs = cfg.Sentinels
		opts.MasterName = cfg.Master
		opts.DB = cfg.DB
	case "credis":
		// credis 兼容普通 redis，直接用 redis 客户端即可
	}
	logger.DebugLog("getRedisClient: mode=%s, addrs=%v, master=%s, db=%d", cfg.Mode, opts.Addrs, opts.MasterName, opts.DB)
	return redis.NewUniversalClient(opts), nil
}

// 连通性测试
func CacheConnect(cfg CacheConfig) map[string]string {
	result := map[string]string{"success": "false"}
	client, err := getRedisClient(cfg)
	if err != nil {
		result["error"] = fmt.Sprintf("client error: %v", err)
		logger.DebugLog("client error: %v", err)
		return result
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()
	logger.DebugLog("CacheConnect: PING")
	_, err = client.Ping(ctx).Result()
	if err != nil {
		result["error"] = fmt.Sprintf("ping error: %v", err)
		logger.DebugLog("ping error: %v", err)
		return result
	}
	result["success"] = "true"
	logger.DebugLog("ping success")
	return result
}

// 写入测试
func CacheWrite(cfg CacheConfig, key, value string) map[string]string {
	result := map[string]string{"success": "false"}
	client, err := getRedisClient(cfg)
	if err != nil {
		result["error"] = fmt.Sprintf("client error: %v", err)
		logger.DebugLog("client error: %v", err)
		return result
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()
	logger.DebugLog("CacheWrite: SET %s %s", key, value)
	err = client.Set(ctx, key, value, 60*time.Second).Err()
	if err != nil {
		result["error"] = fmt.Sprintf("set error: %v", err)
		logger.DebugLog("set error: %v", err)
		return result
	}
	result["success"] = "true"
	logger.DebugLog("set success")
	return result
}

// 删除测试
func CacheDelete(cfg CacheConfig, key string) map[string]string {
	result := map[string]string{"success": "false"}
	client, err := getRedisClient(cfg)
	if err != nil {
		result["error"] = fmt.Sprintf("client error: %v", err)
		logger.DebugLog("client error: %v", err)
		return result
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()
	logger.DebugLog("CacheDelete: DEL %s", key)
	_, err = client.Del(ctx, key).Result()
	if err != nil {
		result["error"] = fmt.Sprintf("del error: %v", err)
		logger.DebugLog("del error: %v", err)
		return result
	}
	result["success"] = "true"
	logger.DebugLog("del success")
	return result
}

// 一键检测
func VerifyCache(cfg CacheConfig) CacheResult {
	key := "precheck_test_key"
	value := "ok"
	res := CacheResult{
		Connect: CacheConnect(cfg),
		Write:   map[string]string{"success": "skip"},
		Delete:  map[string]string{"success": "skip"},
	}
	if res.Connect["success"] == "true" {
		res.Write = CacheWrite(cfg, key, value)
		res.Delete = CacheDelete(cfg, key)
	}
	return res
}

func VerifyCacheJson(cfg CacheConfig) []byte {
	res := VerifyCache(cfg)
	b, _ := json.Marshal(res)
	return b
}
