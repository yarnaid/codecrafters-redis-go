package storage

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
	"time"

	"my-redis/app/parser"
)

type StreamId struct {
	Time, Seq int
}

func (s *StreamId) String() string {
	if s.Time < 0 {
		return "*"
	}
	if s.Seq < 0 {
		return fmt.Sprintf("%d-*", s.Time)
	}
	return fmt.Sprintf("%d-%d", s.Time, s.Seq)
}

type InvalidStreamIdSeq struct {
	Id1, Id2 *StreamId
}

type ZeroStreamId struct{}

var (
	_ error = (*InvalidStreamIdSeq)(nil)
	_ error = (*ZeroStreamId)(nil)
)

func (z ZeroStreamId) Error() string {
	return "stream id must not be 0-0"
}

func (i InvalidStreamIdSeq) Error() string {
	return fmt.Sprintf("wrong streamId seq: %v -> %v", i.Id1, i.Id2)
}

func StreamIdGreaterOrGen(id1, id2 *StreamId, eq bool) bool {
	if id2.Time < 0 {
		id2.Time = int(time.Now().UnixMilli())
		if id2.Time == id1.Time {
			id2.Seq = id1.Seq + 1
		} else {
			id2.Seq = 0
		}
		return false
	}
	if id2.Seq < 0 {
		if id2.Time == id1.Time {
			id2.Seq = id1.Seq + 1
		} else {
			id2.Seq = 0
		}
		return false
	}
	switch c := cmp.Compare(id1.Time, id2.Time); {
	case c < 0:
		return false
	case c > 0:
		return true
	default:
		switch cc := cmp.Compare(id1.Seq, id2.Seq); {
		case cc > 0:
			return true
		case cc == 0:
			return eq
		default:
			return false
		}
	}
}

func (s *StreamId) IsZero() bool {
	if s.Time == 0 && s.Seq == 0 {
		return true
	}
	return false
}

type StreamValue struct {
	K string
	V string
}

type StreamContainer struct {
	Id     StreamId
	Values []StreamValue
}

func (s *StreamContainer) GreaterOrGen(other *StreamContainer, eq bool) bool {
	return StreamIdGreaterOrGen(&s.Id, &other.Id, eq)
}

var _ parser.Serializable = (*StreamContainer)(nil)

func (s *StreamContainer) Serialize() ([]byte, error) {
	res := make(parser.Array[parser.Serializable], len(s.Values)+1)
	res[0] = parser.BulkString(s.Id.String())
	for i, v := range s.Values {
		res[i+1] = parser.Array[parser.Serializable]{parser.BulkString(v.K), parser.BulkString(v.V)}
	}
	return res.Serialize()
}

func NewStreamContainer() *StreamContainer {
	return &StreamContainer{Values: make([]StreamValue, 0)}
}

func ParseStreamId(id string) (StreamId, error) {
	timeS, seqS, found := strings.Cut(id, "-")
	if !found {
		if id == "*" {
			return StreamId{-1, -1}, nil
		}
		return StreamId{}, fmt.Errorf("`-` not found in stream id: %q", id)
	}
	time, err := strconv.Atoi(timeS)
	if err != nil {
		return StreamId{}, err
	}
	seq, err := strconv.Atoi(seqS)
	if err != nil {
		if seqS == "*" {
			return StreamId{time, -1}, nil
		}
		return StreamId{}, err
	}
	if seq == 0 && time == 0 {
		return StreamId{}, ZeroStreamId{}
	}
	return StreamId{time, seq}, nil
}

func ParseRangeId(id string) (StreamId, error) {
	if !strings.Contains(id, "-") {
		time, err := strconv.Atoi(id)
		if err != nil {
			return StreamId{}, err
		}
		return StreamId{time, 0}, nil
	}
	timeS, seqS, _ := strings.Cut(id, "-")
	time, err := strconv.Atoi(timeS)
	if err != nil {
		return StreamId{}, err
	}
	seq, err := strconv.Atoi(seqS)
	if err != nil {
		return StreamId{}, err
	}
	return StreamId{time, seq}, nil
}

func StreamRange(streams []*StreamContainer, startId, endId string) ([]*StreamContainer, error) {
	start, err := ParseRangeId(startId)
	if err != nil {
		return nil, err
	}
	end, err := ParseRangeId(endId)
	if err != nil {
		return nil, err
	}
	startRes := getStreamStart(start, streams)

	endRes := getStreamEnd(end, streams)

	return streams[startRes : endRes+1], nil
}

func getStreamEnd(end StreamId, streams []*StreamContainer) int {
	endContainer := StreamContainer{Id: end}
	for i := len(streams) - 1; i >= 0; i-- {
		s := streams[i]
		if endContainer.GreaterOrGen(s, true) {
			return i
		}
	}
	return 0
}

func getStreamStart(start StreamId, streams []*StreamContainer) int {
	startContainer := StreamContainer{Id: start}
	for i, s := range streams {
		if s.GreaterOrGen(&startContainer, false) {
			return max(0, i-1)
		}
	}
	return len(streams) - 1
}
