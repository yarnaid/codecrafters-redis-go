package storage

func (c *Coordinator) WaitersArrayLen(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if q, ok := c.blpopWaiters[key]; ok {
		return q.Len()
	}
	return 0
}

func (c *Coordinator) WaitersStreamLen(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if q, ok := c.xreadWaiters[key]; ok {
		return q.Len()
	}
	return 0
}

func (c *Coordinator) Backend() Backend {
	return c.backend
}

func EntriesAfter(stream []*StreamContainer, startId StreamId) []*StreamContainer {
	return entriesAfter(stream, startId)
}

func PointersToValues(values []*StreamContainer) []StreamContainer {
	res := make([]StreamContainer, len(values))
	for i := range len(values) {
		res[i] = *values[i]
	}
	return res
}
