package main

import (
	"net/http"

	"github.com/saryginrodion/stackable"
	"github.com/saryginrodion/stackable/middleware"
	"github.com/sirupsen/logrus"
)

type Shared struct{}

type Local struct{}

func (s Local) Default() any {
	return Local{}
}

type Context = stackable.Context[Shared, Local]

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	stack := stackable.NewStackable[Shared, Local](
		new(Shared),
	)

	stack.AddHandler(&middleware.CORSMiddleware[Shared, Local]{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("GET /", stack.AddUniqueHandler(
		stackable.WrapFunc(func(context *Context, next func() error) error {
			context.Response = stackable.NewHttpResponse(
				http.StatusOK,
				"text/html",
				"<h1>Index route!</h1>",
			)
			return next()
		}),
	))

	http.Handle("POST /", stack.AddUniqueHandler(
		stackable.WrapFunc(func(context *Context, next func() error) error {
			context.Response = stackable.NewHttpResponse(
				http.StatusOK,
				"text/html",
				"<h1>Index route (POST method)!</h1>",
			)
			return next()
		}),
	))

	logrus.Info("Starting...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		logrus.Fatal(err)
	}
}
