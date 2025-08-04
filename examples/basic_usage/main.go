package main

import (
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/radyshenkya/stackable"
	"github.com/sirupsen/logrus"
)

// Every handler in handlers stack has access to shared state - there you can store your DB connections and stuff like this.
type Shared struct {
	requestCounter atomic.Int32
}

// LocalState struct instantiated for every request - you can store unique to request data (e. g. user struct got from auth layer)
type Local struct {
	requestId int32
}

// LocalState needs to implement Default interface - with this values it will be instantiated for every request.
func (s Local) Default() any {
	return Local{
		// This field will be changed in requestId handler
		requestId: 0,
	}
}

// For readability, you can declare Context type.
type Context = stackable.Context[Shared, Local]

// Handlers (or layers) in stack is values implementing stackable.Handler[S, L] interface
// SetRequestIdMiddleware will be function Handler - it does not need to store anything, everything for it already saved in Shared.
var SetRequestIdMiddleware = stackable.WrapFunc(
	// context - holds Requst, Response, Shared and Local
	// next - function to call next handler in stack
	func(context *Context, next func() error) error {
		context.Local.requestId = context.Shared.requestCounter.Load()
		context.Shared.requestCounter.Add(1)

		// If you want to call next handler - use return next().
		// If no, you can return nil (handler succeed) or error.
		// When last handler in stack calls next it will return nil.
		return next()
	},
)

// Example of struct middleware
type LoggingMiddleware struct {
	tag string
}

func (s *LoggingMiddleware) Run(context *Context, next func() error) error {
	// Calling next() first to apply every layer below
	err := next()

	logrus.WithFields(
		logrus.Fields{
			"tag": s.tag, // Layers also can store some values inside them, like Shared, but only for the Handler instance
			"rid": context.Local.requestId,
			"ip":  context.Request.RemoteAddr,
		},
	).Infof("%d - %s %s", context.Response.Status(), context.Request.Method, context.Request.URL.Path)

	return err
}

// Layer for mapping errors to Json objects.
var ErrorMapperMiddleware = stackable.WrapFunc(
	func(context *Context, next func() error) error {
		err := next()
		if err != nil {
			var httpErr stackable.HttpError
			if errors.As(err, &httpErr) {
			} else {
				httpErr = stackable.HttpError{
					Status:  http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
			context.Response, _ = stackable.JsonResponse(
				httpErr.Status,
				httpErr,
			)
			return nil
		}
		return err
	},
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	// Creating new Stackable with Shared instance
	stack := stackable.NewStackable[Shared, Local](
		new(Shared),
	)

	stack.SetLogLevel(logrus.DebugLevel)

	// Handlers run from the first added to last
	// AddHandler adds handler to existing Stackable
	// With AddHandler you can add some layers, that will be applied for every request (if you are using this stack)
	stack.AddHandler(SetRequestIdMiddleware)
	stack.AddHandler(ErrorMapperMiddleware)

	// AddUniqueHandler is copying Stackable instance and adds new Handler to this. Stackable in stack will not be touched
	http.Handle("GET /", stack.AddUniqueHandler(&LoggingMiddleware{tag: "index route"}).AddUniqueHandler(
		stackable.WrapFunc(func(context *Context, next func() error) error {
			// Writing response
			context.Response = stackable.NewHttpResponse(
				http.StatusOK,
				"text/html",
				"<h1>Index route!</h1>",
			)
			return next()
		}),
	))

	// This handler will not use LoggingMiddleware.
	http.Handle("GET /json", stack.AddUniqueHandler(
		stackable.WrapFunc(func(context *Context, next func() error) error {
			context.Response, _ = stackable.JsonResponse(
				http.StatusOK,
				struct {
					Message string `json:"msg"`
				}{
					Message: "Hello World!",
				},
			)
			return next()
		}),
	))

	// We can throw errors!
	http.Handle("GET /error", stack.AddUniqueHandler(&LoggingMiddleware{tag: "error route"}).AddUniqueHandler(
		stackable.WrapFunc(func(context *Context, next func() error) error {
			return stackable.HttpError{
				Status:  http.StatusTeapot,
				Message: "I AM A TEAPOT",
			}
		}),
	))

	logrus.Info("Starting...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		logrus.Fatal(err)
	}
}
