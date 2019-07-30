package proxy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultListenPort is the default proxy listen port
	DefaultListenPort = 8888
	// DefaultGlobalCacheExpiry is the default global cache expiry time duration in seconds
	DefaultGlobalCacheExpiry = 300 // 5min
	// DefaultCacheMaxEntry is the default maximum number of cache entries
	DefaultCacheMaxEntry = 200
	// DefaultMaxConcurrentConn is the default maximum number of concurrent connections
	DefaultMaxConcurrentConn = 1000
)

// Config represents the proxy configuration
type Config struct {
	RedisAddr  string
	ListenPort int
	MaxConn    int
}

// server represents an web server that runs the proxy
type server struct {
	redisdb *redis.Client
	config  Config
	logger  *logrus.Logger
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

	return &server{
		redisdb: client,
		config:  config,
		logger:  logger,
	}
}

// API - Method: GET - Endpoint - /get?key=<key>
func (s *server) getKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	// read from redis
	val, err := s.getKeyFromRedis(key)
	if err == redis.Nil {
		http.Error(w, fmt.Sprintf("Entry with key %s doesn't exist", key), http.StatusNotFound)
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		fmt.Fprint(w, val)
	}
}

// getKeyFromRedis performs a GET command with a key to redis
func (s *server) getKeyFromRedis(key string) (string, error) {
	val, err := s.redisdb.Get(key).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		s.logger.Error("failed to retrieve key ", key, "from redis - error: ", err.Error())
		return "", err
	} else {
		return val, nil
	}
}

// Run starts the proxy web server
func Run(config Config) {
	server := newServer(config)
	http.HandleFunc("/get", server.getKey)

	addr := fmt.Sprintf(":%d", server.config.ListenPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
