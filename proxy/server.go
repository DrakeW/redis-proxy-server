package proxy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
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
}

// newServer returns a new server object based on the proxy config input
func newServer(config Config) *server {
	if config.RedisAddr == "" {
		log.Fatal("Address of redis instance cannot be empty")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
		OnConnect: func(con *redis.Conn) error {
			log.Print("Connection to redis established at ", config.RedisAddr)
			return nil
		},
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("Failed to connect to redis at ", config.RedisAddr)
	}
	return &server{
		redisdb: client,
		config:  config,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world at %s", r.URL.Path[1:])
}

// Run starts the proxy web server
func Run(config Config) {
	server := newServer(config)
	http.HandleFunc("/", handler)

	addr := fmt.Sprintf(":%d", server.config.ListenPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
