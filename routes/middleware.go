package routes

import (
	"net/http"
)

type Middleware struct {
	serveHTTP    func(http.ResponseWriter, *http.Request) bool
	prev         *Middleware
	isShortcuted bool
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m.prev != nil {
		m.prev.ServeHTTP(w, req)
		if m.prev.isShortcuted {
			return
		}
	}

	m.isShortcuted = m.serveHTTP(w, req)
}

func (m *Middleware) Then(next Middleware) Middleware {
	next.prev = m
	return next
}

func (m *Middleware) ThenHandler(next http.Handler) http.Handler {
	mw := newMiddlewareHandler(next)
	mw.prev = m
	return mw
}

func (m *Middleware) ThenFunc(next func(http.ResponseWriter, *http.Request)) http.Handler {
	mw := newMiddlewareHandler(http.HandlerFunc(next))
	mw.prev = m
	return mw
}

func newMiddleware(f func(http.ResponseWriter, *http.Request) bool) *Middleware {
	return &Middleware{serveHTTP: f, isShortcuted: true}
}

func newMiddlewareHandler(next http.Handler) *Middleware {
	return &Middleware{serveHTTP: func(w http.ResponseWriter, req *http.Request) bool {
		next.ServeHTTP(w, req)
		return true
	}, isShortcuted: true}
}
