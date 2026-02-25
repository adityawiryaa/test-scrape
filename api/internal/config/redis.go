package config

import "strconv"

type RedisConfig struct {
	Host    string
	Port    string
	DB      int
	AsynqDB int
}

func LoadRedisConfig() *RedisConfig {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	asynqDB, _ := strconv.Atoi(getEnv("ASYNQ_DB", "1"))

	return &RedisConfig{
		Host:    getEnv("REDIS_HOST", "localhost"),
		Port:    getEnv("REDIS_PORT", "6379"),
		DB:      db,
		AsynqDB: asynqDB,
	}
}

func (r *RedisConfig) Addr() string {
	return r.Host + ":" + r.Port
}
