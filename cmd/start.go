/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxpixel/api"
	_ "github.com/danesparza/fxpixel/docs" // swagger docs location
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

	//	Create an api service object
	apiService := api.Service{
		PlayTimeline:     make(chan leds.PlayTimelineRequest),
		StopTimeline:     make(chan string),
		StopAllTimelines: make(chan bool),
		//DB:           appdata,
		StartTime: time.Now(),
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
		//	File management
		r.Put("/audio", apiService.ShowUI)
		r.Get("/audio", apiService.ShowUI)
		r.Delete("/audio/{id}", apiService.ShowUI)
		r.Post("/audio/{id}", apiService.ShowUI)

		//	Play audio
		r.Post("/audio/play", apiService.ShowUI)
		r.Post("/audio/play/{id}", apiService.ShowUI)
		r.Post("/audio/play/random/{tag}", apiService.ShowUI)
		r.Post("/audio/stream", apiService.ShowUI)
		r.Post("/audio/loop/{id}/{loopTimes}", apiService.ShowUI)

		//	Stop audio
		r.Post("/audio/stop", apiService.ShowUI)
		r.Post("/audio/stop/{id}", apiService.ShowUI)

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
