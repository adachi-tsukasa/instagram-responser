package main

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func newContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func logf(c context.Context, format string, args ...interface{}) {
	log.Infof(c, format, args...)
}

func errorf(c context.Context, format string, args ...interface{}) {
	log.Errorf(c, format, args...)
}
