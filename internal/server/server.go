package server

import (
	"net"

	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/repository"
	"go.uber.org/zap"
)

type HandlerFunc func(conn net.Conn, repo *repository.GlobalRepository, logger *zap.Logger)

type Router interface {
	GetRoute(handlerID byte) (HandlerFunc, bool)
}

type Server interface {
	Start() error
	AddHandler(name string, handler HandlerFunc)
	Shutdown() error
}
