package httpbreaker

// MemoryCounts is a memory-based implementation of Counts.
type MemoryCounts struct {
	requests             uint32
	totalSuccesses       uint32
	totalFailures        uint32
	consecutiveSuccesses uint32
	consecutiveFailures  uint32
}

// NewMemoryCounts returns a new instance of MemoryCounts.
func NewMemoryCounts() *MemoryCounts {
	return &MemoryCounts{}
}

// OnRequest increments the request count.
func (c *MemoryCounts) OnRequest() {
	c.requests++
}

// OnSuccess increments the total success and consecutive success counts, and
// resets the consecutive failure count.
func (c *MemoryCounts) OnSuccess() {
	c.totalSuccesses++
	c.consecutiveSuccesses++
	c.consecutiveFailures = 0
}

// OnFailure increments the total failure and consecutive failure counts, and
// resets the consecutive success count.
func (c *MemoryCounts) OnFailure() {
	c.totalFailures++
	c.consecutiveFailures++
	c.consecutiveSuccesses = 0
}

// Clear resets all counts to zero.
func (c *MemoryCounts) Clear() {
	c.requests = 0
	c.totalSuccesses = 0
	c.totalFailures = 0
	c.consecutiveSuccesses = 0
	c.consecutiveFailures = 0
}

// Requests returns the total number of requests.
func (c *MemoryCounts) Requests() uint32 {
	return c.requests
}

// ConsecutiveFailures returns the number of consecutive failures.
func (c *MemoryCounts) ConsecutiveFailures() uint32 {
	return c.consecutiveFailures
}

// ConsecutiveSuccesses returns the number of consecutive successes.
func (c *MemoryCounts) ConsecutiveSuccesses() uint32 {
	return c.consecutiveSuccesses
}

// TotalFailures returns the total number of failures.
func (c *MemoryCounts) TotalFailures() uint32 {
	return c.totalFailures
}

// Counts is an interface for tracking counts of requests and responses.
type Counts interface {
	// OnRequest increments the request count.
	OnRequest()
	// OnSuccess increments the total success and consecutive success counts,
	OnSuccess()
	// OnFailure increments the total failure and consecutive failure counts,
	OnFailure()
	// Clear resets all counts to zero.
	Clear()
	// Requests returns the total number of requests.
	Requests() uint32
	// ConsecutiveFailures returns the number of consecutive failures.
	ConsecutiveFailures() uint32
	// ConsecutiveSuccesses returns the number of consecutive successes.
	ConsecutiveSuccesses() uint32
	// TotalFailures returns the total number of failures.
	TotalFailures() uint32
}
