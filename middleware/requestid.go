package middleware

import (
	"github.com/saryginrodion/stackable"
	"sync/atomic"
)

type ILocalRequestId interface {
	RequestId() int64
	SetRequestId(int64)
}

type LocalRequestId struct {
	requestId int64
}

func (s *LocalRequestId) RequestId() int64 {
	return s.requestId
}

func (s *LocalRequestId) SetRequestId(newId int64) {
	s.requestId = newId
}

type RequestIdMiddleware[S any, L interface{ ILocalRequestId; stackable.Default}] struct {
	counter atomic.Int64
}

func (s *RequestIdMiddleware[S, L]) Run(context *stackable.Context[S, L], next func() error) error {
	(*context.Local).SetRequestId(s.counter.Load())
	s.counter.Add(1)
	return next()
}
