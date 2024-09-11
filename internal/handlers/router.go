package handlers

import (
	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/server"
)

const (
	HandlerQuote byte = 0x01 // Constant for the "quote" handler
	// You can define more handler constants here
)

type Router struct {
	routes map[byte]server.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[byte]server.HandlerFunc),
	}
}

func (r *Router) AddRoute(handlerID byte, handler server.HandlerFunc) {
	r.routes[handlerID] = handler
}

func (r *Router) GetRoute(handlerID byte) (server.HandlerFunc, bool) {
	handler, exists := r.routes[handlerID]
	return handler, exists
}
