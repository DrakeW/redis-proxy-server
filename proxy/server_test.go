package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var proxy = newServer(Config{
	ListenPort:      "8888",
	RedisAddr:       "localhost:6379",
	MaxConn:         1000,
	CacheExpiry:     20,
	CacheMaxEntries: 10,
})

var testCacheKV = []string{"a", "1"}
var testRedisKV = []string{"c", "3"}

func TestGetHandler(t *testing.T) {
	// set up test server
	handler := http.HandlerFunc(proxy.getKey)
	// set up cache data
	proxy.cache.Add(testCacheKV[0], testCacheKV[1])
	// set up redis data
	proxy.redisdb.Set(testRedisKV[0], testRedisKV[1], 20*time.Second).Result()

	t.Run("Success from cache- 200", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest("GET", fmt.Sprintf("/get?key=%s", testCacheKV[0]), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("wrong status code returned: got %v expected %v",
				status, http.StatusOK)
		}
		expected := testCacheKV[1]
		if rr.Body.String() != expected {
			t.Errorf("wrong body returned: got %v expected %v",
				rr.Body.String(), expected)
		}
	})

	t.Run("Success from redis - 200", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest("GET", fmt.Sprintf("/get?key=%s", testRedisKV[0]), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("wrong status code returned: got %v expected %v",
				status, http.StatusOK)
		}
		expected := testRedisKV[1]
		if rr.Body.String() != expected {
			t.Errorf("wrong body returned: got %v expected %v",
				rr.Body.String(), expected)
		}
	})

	t.Run("Not Exist - 404", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest("GET", "/get?key=non-existent", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("wrong status code returned: got %v expected %v",
				status, http.StatusOK)
		}
	})
}
