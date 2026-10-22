package storage

import (
	"container/list"

	"my-redis/app/parser"
)

type waiterArray struct {
	ch   chan parser.Serializable
	elem *list.Element
}

type waiterStream struct {
	ch      chan StreamContainer
	elem    *list.Element
	StartID StreamID
}

func NewWaiterArray() waiterArray {
	return waiterArray{ch: make(chan parser.Serializable, 1)}
}

func NewWaiterStream(startID StreamID) waiterStream {
	return waiterStream{ch: make(chan StreamContainer, 1), StartID: startID}
}
