package setup

import (
	"io"
	"log"

	"github.com/xFidle/sportradar-intern/server/internal/db"
	"github.com/xFidle/sportradar-intern/server/internal/handlers"
	"github.com/xFidle/sportradar-intern/server/internal/service"
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
	event       *service.EventService
	sport       *service.SportService
	competition *service.CompetitionService
	team        *service.TeamService
	venue       *service.VenueService
}

type transport struct {
	event        *handlers.EventHandler
	sport        *handlers.SportHandler
	team         *handlers.TeamHandler
	eventOptions *handlers.EventOptionsHandler
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
	dbHandle := b.db.Handle()
	b.services.event = service.NewEventService(dbHandle, b.config.fileserverAddr)
	b.services.sport = service.NewSportService(dbHandle)
	b.services.competition = service.NewCompetitionService(dbHandle, b.config.fileserverAddr)
	b.services.team = service.NewTeamService(dbHandle, b.config.fileserverAddr)
	b.services.venue = service.NewVenueService(dbHandle)
}

func (b *bootstrap) initTransport() {
	b.transport.event = handlers.NewEventHandler(b.services.event)
	b.transport.sport = handlers.NewSportHandler(b.services.sport)
	b.transport.team = handlers.NewTeamHandler(b.services.team)
	b.transport.eventOptions = handlers.NewEventOptionsHandler(b.services.competition, b.services.venue)
}

func (b *bootstrap) close() {
	for _, closer := range b.closers {
		if err := closer.Close(); err != nil {
			log.Printf("Connection with service was not shutdown properly: %v", err)
		}
	}
}
