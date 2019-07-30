package cmd

import (
	"os"

	"github.com/DrakeW/redis-cache-proxy/proxy"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "redis-proxy",
	Short:        "redis-proxy is a simple read-through Redis proxy that adds caching capability in front of your redis instance",
	Version:      "0.0.1",
	SilenceUsage: true,
	Run:          startProxyServer,
}

func startProxyServer(cmd *cobra.Command, args []string) {
	proxyConfig := proxy.Config{
		ListenPort: listenPort,
		RedisAddr:  redisAddr,
		MaxConn:    maxConnection,
	}
	proxy.Run(proxyConfig)
}

// port the proxy service should listen to
var listenPort int

// address of the backing redis instance
var redisAddr string

// global cache expiry time duration
var cacheExpiry int

// maximum number of keys in cache
var cacheMaxEntry int

// maximum concurrent client connection
var maxConnection int

func init() {
	rootCmd.Flags().IntVarP(&listenPort, "port", "p", proxy.DefaultListenPort, "The port redis-proxy should listen to")
	rootCmd.Flags().StringVar(&redisAddr, "redis-addr", "", "The address of the backing redis instance")
	rootCmd.Flags().IntVar(&cacheExpiry, "cache-expiry", proxy.DefaultGlobalCacheExpiry, "Global cache expiry time duration (in seconds)")
	rootCmd.Flags().IntVar(&cacheMaxEntry, "cache-max-entry", proxy.DefaultCacheMaxEntry, "Maximum number of keys the cache holds at a time")
	rootCmd.Flags().IntVar(&maxConnection, "max-conn", proxy.DefaultMaxConcurrentConn, "Maximum number of concurrent connections the proxy accepts")
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
