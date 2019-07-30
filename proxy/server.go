package proxy

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
