package server

import (
	"net/http"
	"sync/atomic"

	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      	*database.Queries
	JWTSecret      	string
	PolkaKey       	string
	Logger		*Logger
}

func New(cfg *APIConfig, filepathRoot, port string) *http.Server{
	r := chi.NewRouter()
	
	r.Use(middlewareCreateRequestID)
	r.Use(cfg.InputSanatizer)
	r.Use(cfg.MiddlewareMetricsInc)
	r.Use(cfg.MiddlewareCreateUserEvent)


	r.Route("/api", func (r chi.Router){
		//r.Get("/users", cfg.handlerGetUser)
		r.Post("/users", cfg.handlerAddUser)
		//r.Put("/users", cfg.handlerUpdateUser)
		//r.Delete("/users", cfg.handlerDeleteUser)
		r.Get("/healthz", handlerReadiness)
		r.Post("/login", cfg.handlerLogin)
		r.Get("/refresh", cfg.handlerRefresh)
		r.Get("/revoke", cfg.handlerRevoke)
		r.Get("/reset", cfg.handlerReset)
	})

	/*	
	r.Route("/admin", func(r chi.Router){
		r.Use(cfg.middlewareAdmin)

		r.Route("/events", func(r chi.Router){
			r.Get("/", cfg.handlerGetAllEvents)
			r.Get("/latest", cfg.handlerGetLatestEvents)
			r.Get("/user", cfg.handlerGetUserEvents)
			r.Get("/methode", cfg.handlerGetMethodeEvents)
			r.Get("/reset", cfg.handlerResetEvents)
		})

		r.Route("/users", func(r chi.Router){
			r.Get("/user", cfg.handlerGetUsers)
			r.Route("/{id}", func(r chi.Router){
				r.Get("/", cfg.handlerGetUser)
				r.Post("/", cfg.handlerAddUser)
				r.Put("/", cfg.handlerUpdateUser)
				r.Delete("/", cfg.handlerDeleteUser)
			})
		})

		r.Route("/actions", func(r chi.Router){
			r.Get("/", cfg.handlerListActions)
			r.With(cfg.middlewareActionID).Post("/", cfg.handlerExecuteAction)
			r.Route("/{id}", func(r chi.Router){
				r.Put("/", cfg.handlerAdjustAction)
				r.Delete("/", cfg.handlerRemoveAction)	
			})
		})
	})
	*/

	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/app/", http.StatusPermanentRedirect)
	})

	fs := http.FileServer(http.Dir(filepathRoot))
	r.Handle("/app/*", http.StripPrefix("/app", fs))
	
	return &http.Server{
		Addr:  ":" + port,
		Handler: r,
	}
} 
