package stackable

import (
	"io"
	"net/http"
	"time"

	"github.com/saryginrodion/stackable/utils"
	"github.com/sirupsen/logrus"
)

type ILocalState interface {
    Default
}

type ISharedState any

type Context[S ISharedState, L ILocalState] struct {
	Shared   *S
	Local    *L
	Response Response
	Request  *http.Request
}

type Handler[S ISharedState, L ILocalState] interface {
	Run(context *Context[S, L], next func() error) error
}

type Stackable[S ISharedState, L ILocalState] struct {
	Handlers []Handler[S, L]
	Shared   *S
	logger *logrus.Logger
}

func NewStackable[S ISharedState, L ILocalState](s *S) Stackable[S, L] {
	logger := logrus.New()

	stack := Stackable[S, L]{
		Shared: s,
		logger: logger,
	}

	return stack
}

func (s *Stackable[S, L]) SetLogLevel(level logrus.Level) {
	s.logger.SetLevel(level)
}

func (s *Stackable[S, L]) AddHandler(handler Handler[S, L]) *Stackable[S, L] {
	s.Handlers = append(s.Handlers, handler)
	return s
}

func (s Stackable[S, L]) AddUniqueHandler(handler Handler[S, L]) Stackable[S, L] {
	newStackable := Stackable[S, L]{
		Shared: s.Shared,
		logger: s.logger,
	}

	newStackable.Handlers = make([]Handler[S, L], len(s.Handlers))
	copy(newStackable.Handlers, s.Handlers)
	newStackable.Handlers = append(newStackable.Handlers, handler)

	return newStackable
}

func (s *Stackable[S, L]) HttpHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		local := DefaultValue[L]()
		var context Context[S, L] = Context[S, L]{
			Shared: s.Shared,
			Local:  &local,
			Response: NewHttpResponse(
				500,
				"text/html",
				"Override this response from your handlers!",
			),
			Request: request,
		}

		handlerIndex := 0

		var next func() error

		next = func() error {
			if handlerIndex >= len(s.Handlers) {
				return nil
			}

			layer := s.Handlers[handlerIndex]
			handlerIndex += 1

			return layer.Run(&context, next)
		}

		err := next()

		if err != nil {
			s.logger.WithField("err", err).Error("Stackable: finished with error. Error: " + err.Error())
		}

		// Writing response to http.ResponseWriter
        headers := context.Response.Headers()

		for key, values := range (&headers).Entries() {
            response.Header().Del(key)

            for _, value := range values {
                response.Header().Add(key, value)
            }
		}

		response.WriteHeader(context.Response.Status())

		_, err = io.Copy(response, context.Response.Body())

		if err != nil {
			s.logger.WithField("err", err).Error("Stackable: failed to write response. Error: " + err.Error())
		}
	}
}

func (s Stackable[S, L]) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	s.HttpHandler()(response, request)
	stopTime := time.Now()
	elapsedTime := stopTime.Sub(startTime)

	s.logger.WithFields(logrus.Fields{
		"origin": "stackable",
		"url": request.URL.Path,
	}).Debugf("Handled in %s", utils.FormatDuration(elapsedTime))
}
