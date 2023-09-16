package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxpixel/api"
	_ "github.com/danesparza/fxpixel/docs" // swagger docs location
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/leds"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API service",
	Long:  `Start the API service`,
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Debug().Str("configFile", viper.ConfigFileUsed()).Msg("Using config file")
	} else {
		log.Debug().Msg("No config file found")
	}

	systemdb := viper.GetString("datastore.system")

	//	Emit what we know:
	log.Info().
		Str("systemdb", systemdb).
		Msg("Starting up")

	//	Init SQLite
	db, err := data.InitSqlite(systemdb)
	if err != nil {
		log.Err(err).Msg("Problem trying to open the system database")
		return
	}
	defer db.Close()

	//	Init the AppDataService
	appdata := data.NewAppDataService(db)

	//	Create an api service object
	apiService := api.Service{
		PlayTimeline:     make(chan leds.PlayTimelineRequest),
		StopTimeline:     make(chan string),
		StopAllTimelines: make(chan bool),
		DB:               appdata,
		StartTime:        time.Now(),
	}

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

	//	Create a router and set up our REST endpoints...
	r := chi.NewRouter()

	//	Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(api.ApiVersionMiddleware)

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

		//	Timeline management
		r.Route("/timelines", func(r chi.Router) {
			r.Put("/", apiService.ShowUI)          // Add a timeline
			r.Get("/", apiService.GetAllTimelines) // Get all timelines
			r.Delete("/{id}", apiService.ShowUI)   // Delete a timeline
			r.Post("/{id}", apiService.ShowUI)     // Update a timeline
		})

		//	Run or stop a timeline
		r.Route("/timeline", func(r chi.Router) {
			r.Post("/run", apiService.ShowUI)              // Run a random timeline
			r.Post("/run/{id}", apiService.ShowUI)         // Run a specific timeline
			r.Post("/run/random/{tag}", apiService.ShowUI) // Run a random timeline in a tag
			r.Post("/stop", apiService.ShowUI)             // Stop all running timelines
			r.Post("/stop/{id}", apiService.ShowUI)        // Stop a specific timeline
		})
	})

	//	SWAGGER
	r.Mount("/swagger", httpSwagger.WrapHandler)

	//	Start the media processor:
	//go media.HandleAndProcess(ctx, apiService.PlayMedia, apiService.StopMedia, apiService.StopAllMedia)

	//	Format the bound interface:
	formattedServerPort := fmt.Sprintf(":%v", viper.GetString("server.port"))

	//	Start the service and display how to access it
	log.Info().Str("server", formattedServerPort).Msg("Started REST service")
	log.Err(http.ListenAndServe(formattedServerPort, r)).Msg("HTTP API service error")
}

func handleSignals(ctx context.Context, sigs <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigs:
		switch sig {
		case os.Interrupt:
			log.Info().Msg("SIGINT")
		case syscall.SIGTERM:
			log.Info().Msg("SIGTERM")
		}

		log.Info().Msg("Shutting down ...")
		cancel()
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
