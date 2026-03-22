package setup

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/xFidle/sportradar-intern/server/internal/db"
)

type bootstrap struct {
	config   config
	database *db.Service
	closers  []io.Closer
}

func newBootstrap() *bootstrap {
	conf := loadConfig()
	database := db.Must(conf.database)
	return &bootstrap{
		config:   conf,
		database: database,
		closers:  []io.Closer{database},
	}
}

func (b *bootstrap) close() {
	for _, closer := range b.closers {
		if err := closer.Close(); err != nil {
			log.Printf("Connection with service was not shutdown properly: %v", err)
		}
	}
}

type App struct {
	Server *http.Server
	Close  func()
}

func NewApp() *App {
	b := newBootstrap()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", b.config.port),
		Handler:      registerRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &App{Server: server, Close: b.close}
}
