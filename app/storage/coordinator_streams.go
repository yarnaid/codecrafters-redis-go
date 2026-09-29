package storage

import (
	"container/list"
	"fmt"
	"log/slog"
	"time"
)

func (c *Coordinator) Xadd(key, id string, kvPairs ...string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.backend.Get(key)
	if !ok {
		v = Value{Kind: KindStream, Value: make([]*StreamContainer, 0)}
	}
	stream, ok := v.Value.([]*StreamContainer)
	if !ok {
		return id, fmt.Errorf("value %T is not slice stream containers k=%q, v=%v", v.Value, key, v.Value)
	}
	newId, err := ParseStreamId(id)
	if err != nil {
		return id, err
	}
	data := make([]StreamValue, int(len(kvPairs)/2))
	for i := 0; i < len(kvPairs); i += 2 {
		data[i] = StreamValue{K: kvPairs[i], V: kvPairs[i+1]}
	}
	newStreamVal := &StreamContainer{Id: newId, Values: data}
	if len(stream) > 0 {
		last := stream[len(stream)-1]
		if last.GreaterOrGen(newStreamVal, false) || (last.Id.Seq == newStreamVal.Id.Seq && last.Id.Time == newStreamVal.Id.Time) {
			err = InvalidStreamIdSeq{&last.Id, &newStreamVal.Id}
			slog.Error("invalid stream id seq", "err", err)
			return newStreamVal.Id.String(), err
		}
		slog.Debug("[Coordinator][Xadd] greater!", "key", key, "oldId", last.Id, "newId", newId)
	} else {
		slog.Debug("[Coordinator][Xadd] skipping empty list validation", "key", key)
		if newStreamVal.Id.Time < 0 {
			newStreamVal.Id = StreamId{int(time.Now().UnixMilli()), 0}
		} else {
			if newStreamVal.Id.Seq < 0 {
				if newStreamVal.Id.Time == 0 {
					newStreamVal.Id.Seq = 1
				} else {
					newStreamVal.Id.Seq = 0
				}
			}
		}
	}

	v.Value = append(stream, newStreamVal)
	c.backend.Set(key, &v)
	c.broadcastStreamValue(key, *newStreamVal)

	return newStreamVal.Id.String(), nil
}

func (c *Coordinator) broadcastStreamValue(key string, value StreamContainer) {
	waiters, ok := c.xreadWaiters[key]
	if !ok {
		return
	}

	for el := waiters.Front(); el != nil; {
		next := el.Next()
		w := el.Value.(waiterStream)
		if StreamIdGreaterOrGen(&value.Id, &w.StartId, false) {
			waiters.Remove(el)
			w.ch <- value
		}
		el = next
	}
	if waiters.Len() == 0 {
		delete(c.xreadWaiters, key)
	}
}

func (c *Coordinator) XRange(key, startId, endId string) ([]*StreamContainer, error) {
	streams, err := c.streamLocked(key)
	if err != nil {
		return nil, err
	}

	return StreamRange(streams, startId, endId)
}

func (c *Coordinator) waitForStreamResult(key string, timeout time.Duration, excludeId StreamId) ([]*StreamContainer, error) {
	waiter := NewWaiterStream(excludeId)
	waiters, ok := c.xreadWaiters[key]
	if !ok {
		waiters = list.New()
		c.xreadWaiters[key] = waiters
	}
	waiters.PushBack(waiter)
	waiter.elem = waiters.Back()
	c.mu.Unlock()

	var timer_ch <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		timer_ch = timer.C
	}

	select {
	case val := <-waiter.ch:
		return []*StreamContainer{&val}, nil
	case <-timer_ch:
		c.mu.Lock()
		defer c.mu.Unlock()
		select {
		case val := <-waiter.ch:
			return []*StreamContainer{&val}, nil
		default:
		}
		waiters := c.xreadWaiters[key]
		waiters.Remove(waiter.elem)
		if waiters.Len() == 0 {
			delete(c.xreadWaiters, key)
		}
		return nil, NewTimeoutError(key, timeout, &excludeId)
	}
}

func (c *Coordinator) XReadStreams(key, id string, timeout *time.Duration) ([]*StreamContainer, error) {
	slog.Debug("[Coordinator][XReadStreams]", "key", key, "id", id, "timeout", deref(timeout))
	start, err := ParseStreamId(id)
	if err != nil {
		switch err.(type) {
		case ZeroStreamId:
		default:
			return nil, err
		}
	}

	c.mu.Lock()
	stream, err := c.streamLocked(key)
	if start.Time == -1 {
		if len(stream) == 0 {
			start = StreamId{0, 0}
		} else {
			last := stream[len(stream)-1]
			start = last.Id
		}
	}
	if err != nil {
		c.mu.Unlock()
		return nil, err
	}
	if res := entriesAfter(stream, start); len(res) > 0 {
		c.mu.Unlock()
		return res, nil
	}
	return c.waitForStreamResult(key, *timeout, start)
}

// streamLocked returns the stream for key. Must be called with c.mu held.
func (c *Coordinator) streamLocked(key string) ([]*StreamContainer, error) {
	v, ok := c.backend.Get(key)
	if !ok {
		return nil, nil
	}
	stream, ok := v.Value.([]*StreamContainer)
	if !ok {
		return nil, &WrongTypeError{v.Value, []*StreamContainer{}}
	}
	return stream, nil
}
