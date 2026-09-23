package storage

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"my-redis/app/parser"
)

type StreamId struct {
	Time, Seq int
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

func StreamIdGreater(id1, id2 StreamId) bool {
	switch c := cmp.Compare(id1.Time, id2.Time); {
	case c < 0:
		return false
	case c > 0:
		return true
	default:
		switch cc := cmp.Compare(id1.Seq, id2.Seq); {
		case cc > 0:
			return true
		default:
			return false
		}
	}
}

func (s StreamId) IsZero() bool {
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

func (s *StreamContainer) Greater(other *StreamContainer) bool {
	return StreamIdGreater(s.Id, other.Id)
}

var _ parser.Serializable = (*StreamContainer)(nil)

func (s *StreamContainer) Serialize() ([]byte, error) {
	return nil, nil
}

func NewStreamContainer() *StreamContainer {
	return &StreamContainer{Values: make([]StreamValue, 0)}
}

func ParseStreamId(id string) (StreamId, error) {
	timeS, seqS, found := strings.Cut(id, "-")
	if !found {
		return StreamId{}, fmt.Errorf("`-` not found in stream id: %q", id)
	}
	time, err := strconv.Atoi(timeS)
	if err != nil {
		return StreamId{}, err
	}
	seq, err := strconv.Atoi(seqS)
	if err != nil {
		return StreamId{}, err
	}
	if seq == 0 && time == 0 {
		return StreamId{}, ZeroStreamId{}
	}
	return StreamId{time, seq}, nil
}
