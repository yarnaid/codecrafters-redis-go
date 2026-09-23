package storage

func (c *Coordinator) WaitersLen(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if q, ok := c.blpop_waiters[key]; ok {
		return q.Len()
	}
	return 0
}

func (c *Coordinator) Backend() Backend {
	return c.backend
}
