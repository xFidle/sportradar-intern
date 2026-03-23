package setup

import (
	"io"
	"log"

	"github.com/xFidle/sportradar-intern/server/internal/db"
	"github.com/xFidle/sportradar-intern/server/internal/event"
	"github.com/xFidle/sportradar-intern/server/internal/handlers"
)

type bootstrap struct {
	config  config
	closers []io.Closer
	storage
	services
	transport
}

type storage struct {
	db *db.Service
}

type services struct {
	event *event.Service
}

type transport struct {
	event *handlers.EventHandler
}

func newBootstrap() *bootstrap {
	conf := loadConfig()
	b := &bootstrap{config: conf, closers: make([]io.Closer, 0)}

	b.initStorage()
	b.initServices()
	b.initTransport()

	return b
}

func (b *bootstrap) initStorage() {
	b.db = db.Must(b.config.database)
	b.closers = append(b.closers, b.db)
}

func (b *bootstrap) initServices() {
	b.services.event = event.New(b.db.Handle(), b.config.fileserverAddr)
}

func (b *bootstrap) initTransport() {
	b.transport.event = handlers.NewEventHandler(b.services.event)
}

func (b *bootstrap) close() {
	for _, closer := range b.closers {
		if err := closer.Close(); err != nil {
			log.Printf("Connection with service was not shutdown properly: %v", err)
		}
	}
}
