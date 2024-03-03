package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Redis struct {
	client *redis.Client // Redis underlying client
	logger *zap.SugaredLogger
	// rdbBkr  breaker.Breaker
	timeout time.Duration
}

// Function creates a new Fascmile Redis source
func NewRedis(addr string, password string, dbIdx int, logger *zap.SugaredLogger, timeout int) *Redis {
	opts := redis.Options{
		Addr:       addr,
		MaxRetries: 2,
		Password:   password,
		DB:         dbIdx,
	}

	return &Redis{
		client: redis.NewClient(&opts),
		logger: logger,
		// rdbBkr:  redisBkr,
		timeout: time.Duration(timeout) * time.Millisecond,
	}
}

var (
	// Errors
	ErrExists = fmt.Errorf("cache entry already exists")
	ErrMiss   = fmt.Errorf("cache miss")
	ErrAdd    = fmt.Errorf("unable to add")
)

// Method inserts the data into the redis
func (r *Redis) Add(ctx context.Context, key string, val []byte, d time.Duration) error {
	st := r.client.WithTimeout(r.timeout).SetNX(ctx, key, val, d)
	if ok, err := st.Result(); err != nil {
		r.logger.Errorw("Redis", "function", "Add", "error", err)
		// r.rdbBkr.LogFailure()
		return err
		// fmt.Printf("%#v", err)
	} else if !ok {
		// r.rdbBkr.LogSuccess()
		// fmt.Printf(`Key does not exist in redis. {"key":%s,"hash":%s}`, key, val)
		return ErrExists
	}
	// r.rdbBkr.LogSuccess()
	// fmt.Printf(`Cache added. {"key":%s,"hash":%s}`, key, val)
	return nil
}

// Method returns the model
func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	// Register
	// gobRegister(h)

	// Get the datain

	st := r.client.WithTimeout(r.timeout).Get(ctx, key)
	value, err := st.Bytes()
	if err == redis.Nil {
		// r.rdbBkr.LogSuccess()
		r.logger.Debugf(`Key doesnt exist in redis. {"key":%s}`, key)
		return nil, ErrMiss
	} else if err != nil {
		// r.rdbBkr.LogFailure()
		r.logger.Errorw("Redis", "function", "Get", "error", err, "key", key)
		return nil, err
	}
	// r.rdbBkr.LogSuccess()

	return value, nil
}

// Method deletes the entry in redis
func (r *Redis) Delete(ctx context.Context, key string) error {
	cmd := r.client.WithTimeout(r.timeout).Del(ctx, key)
	if v, err := cmd.Result(); err != nil {
		// if r.rdbBkr != nil {
		// 	r.rdbBkr.LogFailure()
		// }
		r.logger.Errorw("Redis", "function", "Delete", "error", err, "key", key)
		return err
	} else if v == 0 {
		// if r.rdbBkr != nil {
		// 	r.rdbBkr.LogSuccess()
		// }
		r.logger.Debugf(`Key miss. {"key":%s}`, key)
		return ErrMiss
	}
	// if r.rdbBkr != nil {
	// 	r.rdbBkr.LogSuccess()
	// }
	r.logger.Infof(`Key deleted. {"key":%s}`, key)
	return nil
}

// Method returns a Redis iterator
func (r *Redis) Scan(ctx context.Context, keyPattern string) ([]string, error) {

	keys, _, err := r.client.Scan(ctx, 0, keyPattern, 50).Result()
	if err != nil {
		// r.rdbBkr.LogFailure()
		return nil, err
	}
	// r.rdbBkr.LogSuccess()
	return keys, nil
}
