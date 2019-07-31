package proxy

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DrakeW/redis-cache-proxy/cache"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultListenPort is the default proxy listen port
	DefaultListenPort = "8888"
	// DefaultGlobalCacheExpiry is the default global cache expiry time duration in seconds
	DefaultGlobalCacheExpiry = 300 // 5min
	// DefaultCacheMaxEntry is the default maximum number of cache entries
	DefaultCacheMaxEntry = 200
	// DefaultMaxConcurrentConn is the default maximum number of concurrent connections
	DefaultMaxConcurrentConn = 1000
)

// Config represents the proxy configuration
type Config struct {
	RedisAddr       string
	ListenPort      string
	MaxConn         uint
	CacheExpiry     time.Duration
	CacheMaxEntries uint
}

// server represents an web server that runs the proxy
type server struct {
	redisdb *redis.Client
	config  Config
	logger  *logrus.Logger
	cache   *cache.LRU
}

// newServer returns a new server object based on the proxy config input
func newServer(config Config) *server {
	// initialize logging
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	// set up redis client
	if config.RedisAddr == "" {
		logger.Fatal("address of redis instance cannot be empty")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
		OnConnect: func(con *redis.Conn) error {
			logger.Info("connection to redis established at ", config.RedisAddr)
			return nil
		},
	})
	_, err := client.Ping().Result()
	if err != nil {
		logger.Fatal("failed to connect to redis at ", config.RedisAddr)
	}
	// initialize cache
	cache := cache.NewLRUCache(&cache.Config{
		Expiry:     config.CacheExpiry * time.Second,
		MaxEntries: config.CacheMaxEntries,
	})
	return &server{
		redisdb: client,
		config:  config,
		logger:  logger,
		cache:   cache,
	}
}

// API - Method: GET - Endpoint - /get?key=<key>
func (s *server) getKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	// read from cache
	val := s.cache.Get(key)
	if val == nil {
		// read from redis
		s.getKeyFromRedis(w, key)
	} else {
		valStr, ok := val.(string)
		// if the cached value is not string, it's probably corrupted and should be retrieved using a different cmd
		if !ok {
			s.getKeyFromRedis(w, key)
		} else {
			s.logger.Infof("successfully retrieved key \"%s\" from cache", key)
			fmt.Fprintf(w, valStr)
		}
	}
}

// getKeyFromRedis performs a GET command with a key to redis
func (s *server) getKeyFromRedis(w http.ResponseWriter, key string) {
	val, err := s.redisdb.Get(key).Result()
	if err == redis.Nil {
		http.Error(w, fmt.Sprintf("Entry with key %s doesn't exist", key), http.StatusNotFound)
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.logger.Errorf("failed to retrieve key \"%s\" from redis - error: %s", key, err.Error())
	} else {
		// add retrieved key-value to cache
		s.logger.Infof("successfully retrieved key \"%s\" from redis", key)
		err = s.cache.Add(key, val)
		if err != nil {
			// log the error but still returns the value
			s.logger.Errorf("failed to add key \"%s\" to cache", key)
		}
		fmt.Fprint(w, val)
	}
}

// Run starts the proxy web server
func Run(config Config) {
	server := newServer(config)
	// TODO: add max connection middleware
	http.HandleFunc("/get", server.getKey)

	addr := fmt.Sprintf(":%s", server.config.ListenPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
