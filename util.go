package main

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// newContext は appengine.NewContext を短く書くための関数
func newContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// logf は log.Infof を短く書くための関数
func logf(c context.Context, format string, args ...interface{}) {
	log.Infof(c, format, args...)
}

// errorf は log.Errorf を短く書くための関数
func errorf(c context.Context, format string, args ...interface{}) {
	log.Errorf(c, format, args...)
}
