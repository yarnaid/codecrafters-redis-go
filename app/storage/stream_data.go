package storage

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"my-redis/app/parser"
)

type StreamID struct {
	Time, Seq int
}

func (s *StreamID) String() string {
	if s.Time < 0 {
		return "*"
	}
	if s.Seq < 0 {
		return fmt.Sprintf("%d-*", s.Time)
	}
	return fmt.Sprintf("%d-%d", s.Time, s.Seq)
}

type InvalidStreamIDSeq struct {
	ID1, ID2 *StreamID
}

type ZeroStreamID struct{}

var (
	_ error = (*InvalidStreamIDSeq)(nil)
	_ error = (*ZeroStreamID)(nil)
)

func (z ZeroStreamID) Error() string {
	return "stream id must not be 0-0"
}

func (i InvalidStreamIDSeq) Error() string {
	return fmt.Sprintf("wrong streamId seq: %v -> %v", *i.ID1, *i.ID2)
}

func StreamIDGreaterOrGen(id1, id2 *StreamID, eq bool) bool {
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

func (s *StreamID) IsZero() bool {
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
	ID     StreamID
	Values []StreamValue
}

func (s *StreamContainer) GreaterOrGen(other *StreamContainer, eq bool) bool {
	return StreamIDGreaterOrGen(&s.ID, &other.ID, eq)
}

var _ parser.Serializable = (*StreamContainer)(nil)

func (s *StreamContainer) Serialize() ([]byte, error) {
	res := make(parser.Array[parser.Serializable], len(s.Values)+1)
	res[0] = parser.BulkString(s.ID.String())
	for i, v := range s.Values {
		res[i+1] = parser.Array[parser.Serializable]{parser.BulkString(v.K), parser.BulkString(v.V)}
	}
	return res.Serialize()
}

func NewStreamContainer() *StreamContainer {
	return &StreamContainer{Values: make([]StreamValue, 0)}
}

func ParseStreamID(id string) (StreamID, error) {
	timeS, seqS, found := strings.Cut(id, "-")
	if !found {

		switch id {
		case "*", "$":
			return StreamID{-1, -1}, nil
		}
		return StreamID{}, fmt.Errorf("`-` not found in stream id: %q", id)
	}
	time, err := strconv.Atoi(timeS)
	if err != nil {
		return StreamID{}, err
	}
	seq, err := strconv.Atoi(seqS)
	if err != nil {
		if seqS == "*" {
			return StreamID{time, -1}, nil
		}
		return StreamID{}, err
	}
	if seq == 0 && time == 0 {
		return StreamID{}, ZeroStreamID{}
	}
	return StreamID{time, seq}, nil
}

func ParseRangeID(id string) (StreamID, error) {
	if !strings.Contains(id, "-") {
		time, err := strconv.Atoi(id)
		if err != nil {
			return StreamID{}, err
		}
		return StreamID{time, 0}, nil
	}
	timeS, seqS, _ := strings.Cut(id, "-")
	time, err := strconv.Atoi(timeS)
	if err != nil {
		return StreamID{}, err
	}
	seq, err := strconv.Atoi(seqS)
	if err != nil {
		return StreamID{}, err
	}
	return StreamID{time, seq}, nil
}

func StreamRange(streams []*StreamContainer, startID, endID string) ([]*StreamContainer, error) {
	var start, end StreamID
	var err error
	var startRes, endRes int
	if startID == "-" {
		startRes = 0
	} else {
		start, err = ParseRangeID(startID)
		if err != nil {
			return nil, err
		}
		startRes = getStreamStart(start, streams)
	}
	if endID == "+" {
		endRes = len(streams) - 1
	} else {
		end, err = ParseRangeID(endID)
		if err != nil {
			return nil, err
		}
		endRes = getStreamEnd(end, streams)
	}
	slog.Debug("[StreamRange]", "startId", startID, "endId", endID, "startRes", startRes, "endRes", endRes)
	if startRes < 0 || endRes < 0 {
		return []*StreamContainer{}, nil
	}
	return streams[startRes : endRes+1], nil
}

func getStreamEnd(end StreamID, streams []*StreamContainer) int {
	endContainer := StreamContainer{ID: end}
	for i, s := range slices.Backward(streams) {
		if endContainer.GreaterOrGen(s, true) {
			return i
		}
	}
	return 0
}

func getStreamStart(start StreamID, streams []*StreamContainer) int {
	startContainer := StreamContainer{ID: start}
	for i, s := range streams {
		if s.GreaterOrGen(&startContainer, false) {
			return max(0, i-1)
		}
	}
	return len(streams) - 1
}

func entriesAfter(stream []*StreamContainer, startID StreamID) []*StreamContainer {
	if len(stream) < 1 {
		return stream
	}
	start := getStreamStart(startID, stream)
	firstID := stream[0].ID
	var first int
	if !StreamIDGreaterOrGen(&firstID, &startID, false) {
		first = 1
	}
	return stream[start+first:]
}
