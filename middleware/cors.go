package middleware

import (
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/saryginrodion/stackable"
	"github.com/sirupsen/logrus"
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
            "",
        )

        return err
    }

    bodyBuf := new(strings.Builder)
    _, copyErr := io.Copy(bodyBuf, context.Response.Body())

    if copyErr != nil {
        logrus.Errorln("cors.go: error on copying to bodyBuf: " + copyErr.Error())
    }


    context.Response = stackable.NewHttpResponseRaw(
        headers,
        context.Response.Status(),
        bodyBuf.String(),
    )

    return err
}
