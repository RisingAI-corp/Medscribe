package utils

import (
	"fmt"
	"net/http"
	"sync"
)

// GetAccessToken retrieves the access token from cookies.
func GetAccessToken(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("access token not found")
}

// GetRefreshToken retrieves the refresh token from cookies.
func GetRefreshToken(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("refresh token not found")
}


type SafeResponseWriter struct {
    http.ResponseWriter
    mu sync.Mutex
}
func (w *SafeResponseWriter) Write(p []byte) (n int, err error) {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.ResponseWriter.Write(p)
}

func (w *SafeResponseWriter) Flush() {
    w.mu.Lock()
    defer w.mu.Unlock()
    if f, ok := w.ResponseWriter.(http.Flusher); ok {
        f.Flush()
    }
}

type SafeMap[V any] struct {
    mu sync.Mutex
    m  map[string]V
}

func NewSafeMap[V any]() *SafeMap[V] {
    return &SafeMap[V]{
            m: make(map[string]V),
    }
}

func (sm *SafeMap[V]) Set(key string, value V) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.m[key] = value
}

func (sm *SafeMap[V]) Get(key string) (V, bool) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    val, ok := sm.m[key]
    return val, ok
}

func (sm *SafeMap[V]) Delete(key string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.m, key)
}

func (sm *SafeMap[V]) GetMap() map[string]V {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    // Create a copy of the map
    copyMap := make(map[string]V, len(sm.m))
    for k, v := range sm.m {
            copyMap[k] = v
    }
    return copyMap
}