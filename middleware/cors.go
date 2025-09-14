package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/saryginrodion/stackable"
)


type CORSMiddleware[S stackable.ISharedState, L stackable.ILocalState] struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
    AllowCredentials bool
}

func (s *CORSMiddleware[S, L]) Run(context *stackable.Context[S, L], next func() error) error {
    err := next()

    origin := context.Request.Header.Get("origin")
    headers := context.Response.Headers()

    if slices.Contains(s.AllowedOrigins, "*") {
        headers.Add("Access-Control-Allow-Origin", "*")
    } else if slices.Contains(s.AllowedOrigins, origin) {
        headers.Add("Access-Control-Allow-Origin", origin)
    }

    if slices.Contains(s.AllowedHeaders, "*") {
        headers.Add("Access-Control-Allow-Headers", "*")
    } else {
        headers.Add("Access-Control-Allow-Headers", strings.Join(s.AllowedHeaders, ", "))
    }

    if slices.Contains(s.AllowedMethods, "*") {
        headers.Add("Access-Control-Allow-Methods", "*")
    } else {
        headers.Add("Access-Control-Allow-Methods", strings.Join(s.AllowedMethods, ", "))
    }

    if context.Request.Method == "OPTIONS" {
        context.Response = stackable.NewHttpResponseRaw(
            headers,
            http.StatusOK,
            strings.NewReader(""),
        )

        return err
    }

    bodyStream := context.Response.Body()

    context.Response = stackable.NewHttpResponseRaw(
        headers,
        context.Response.Status(),
        bodyStream,
    )

    return err
}
