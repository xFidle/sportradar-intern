package setup

import (
	"fmt"
	"net/http"
	"time"
)

type App struct {
	Server *http.Server
	Close  func()
}

func NewApp() *App {
	b := newBootstrap()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", b.config.port),
		Handler:      registerRoutes(b.transport),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &App{Server: server, Close: b.close}
}
