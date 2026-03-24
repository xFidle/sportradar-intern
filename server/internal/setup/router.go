package setup

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func registerRoutes(t transport) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "ws://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/event", func(r chi.Router) {
			r.Get("/{event_id}", t.event.HandleGetEvent)
			r.Post("/", t.event.HandleGetEvents) // It is a post method because filter is needed as body
		})

		r.Route("/sports", func(r chi.Router) {
			r.Get("/{sport_id}", t.sport.HandleGetSports)
			r.Get("/{sport_id}/competitions", t.eventOptions.HandleGetCompetitionsBySport)
			r.Get("/{sport_id}/event-options", t.eventOptions.HandleGetEventOptionsBySport)
		})

		r.Route("/competitions", func(r chi.Router) {
			r.Get("/{competition_id}/teams", t.team.HandleGetTeamsByCompetition)
		})
	})

	return r
}
