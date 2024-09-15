package cache

import "time"

// ICache interface with generics
type ICache[T any] interface {
	Get(key interface{}) (T, bool)
	Set(key interface{}, val T) bool
}

// Config struct for Redis configuration
type Config struct {
	Addr     string        `mapstructure:"addr"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	TTL      time.Duration `mapstructure:"ttl"`
}
