package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/go-redis/cache/v9"
)

// CacheCounts is a cache-based implementation of Counts.
type CacheCounts struct {
	client      *cache.Cache
	serviceName string
}

type counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// NewCacheCounts returns a new instance of CacheCounts.
// serviceName is the name of the service.
// client is the cache client.
func NewCacheCounts(client *cache.Cache, serviceName string) *CacheCounts {
	return &CacheCounts{
		client:      client,
		serviceName: serviceName,
	}
}

// OnRequest get count from cache and increment the request count.
func (c *CacheCounts) OnRequest() {
	cts := c.get()
	cts.Requests++
	c.save(cts)
}

// OnSuccess increments the total success and consecutive success counts, and
// resets the consecutive failure count.
func (c *CacheCounts) OnSuccess() {
	cts := c.get()
	cts.TotalSuccesses++
	cts.ConsecutiveSuccesses++
	cts.ConsecutiveFailures = 0
	c.save(cts)
}

// OnFailure increments the total failure and consecutive failure counts, and
// resets the consecutive success count.
func (c *CacheCounts) OnFailure() {
	cts := c.get()
	cts.TotalFailures++
	cts.ConsecutiveFailures++
	cts.ConsecutiveSuccesses = 0
	c.save(cts)
}

// Clear resets all counts to zero.
func (c *CacheCounts) Clear() {
	cts := c.get()
	cts.Requests = 0
	cts.TotalSuccesses = 0
	cts.TotalFailures = 0
	cts.ConsecutiveSuccesses = 0
	cts.ConsecutiveFailures = 0
	c.save(cts)
}

// Requests returns the total number of requests.
func (c *CacheCounts) Requests() uint32 {
	return c.get().Requests
}

// ConsecutiveFailures returns the number of consecutive failures.
func (c *CacheCounts) ConsecutiveFailures() uint32 {
	return c.get().ConsecutiveFailures
}

// ConsecutiveSuccesses returns the number of consecutive successes.
func (c *CacheCounts) ConsecutiveSuccesses() uint32 {
	return c.get().ConsecutiveSuccesses
}

// TotalFailures returns the total number of failures.
func (c *CacheCounts) TotalFailures() uint32 {
	return c.get().TotalFailures
}

func (c *CacheCounts) save(cts counts) {
	err := c.client.Set(&cache.Item{
		Key:   c.getKey(),
		Value: cts,
	})

	if err != nil {
		fmt.Printf("service %s save cache error: %v", c.serviceName, err)
	}
}

func (c *CacheCounts) get() counts {
	var cts counts
	err := c.client.Get(context.Background(), c.getKey(), &cts)

	if err != nil {
		fmt.Printf("service %s get cache error: %v", c.serviceName, err)
	}

	return cts
}

// getkey generate hash by serviceName by md5
func (c *CacheCounts) getKey() string {
	hash := sha256.New()
	hash.Write([]byte(c.serviceName))
	return string(hash.Sum(nil))
}
