package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(apiService Service) http.Handler {
	//	Create a router and set up our REST endpoints...
	r := chi.NewRouter()

	//	Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(ApiVersionMiddleware)

	//	... including CORS middleware
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1", func(r chi.Router) {

		//	System config
		r.Route("/config", func(r chi.Router) {
			r.Get("/", apiService.GetSystemConfig) // Get all system config keys and values
			r.Post("/{key}", apiService.ShowUI)    // Update system config value
		})

		//	Timeline management
		r.Route("/timelines", func(r chi.Router) {
			r.Put("/", apiService.AddTimeline)                     // Add a timeline
			r.Get("/", apiService.GetAllTimelines)                 // Get all timelines
			r.Get("/tag/{tag}", apiService.GetAllTimelinesWithTag) // Get all timelines with a tag
			r.Get("/{id}", apiService.GetTimeline)                 // Get a single timeline
			r.Delete("/{id}", apiService.DeleteTimeline)           // Delete a timeline
			r.Post("/{id}", apiService.UpdateTags)                 // Update timeline tags
		})

		//	Run or stop a timeline
		r.Route("/timeline", func(r chi.Router) {
			//r.Post("/run/random/{tag}", apiService.ShowUI)        // Run a random timeline in a tag
			r.Post("/run/{id}", apiService.RequestTimelinePlay)   // Run a specific timeline
			r.Post("/stop", apiService.RequestAllTimelinesStop)   // Stop all running timeline processes
			r.Post("/stop/{pid}", apiService.RequestTimelineStop) // Stop a specific timeline process
		})
	})

	//	SWAGGER
	r.Mount("/swagger", httpSwagger.WrapHandler)

	return r
}
