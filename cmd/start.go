package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxpixel/api"
	_ "github.com/danesparza/fxpixel/docs" // swagger docs location
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/leds"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	//	Create a background service object
	backgroundService := leds.BackgroundProcess{
		PlayTimeline:     make(chan leds.PlayTimelineRequest),
		StopTimeline:     make(chan string),
		StopAllTimelines: make(chan bool),
		DB:               appdata,
	}

	//	Create an api service object
	apiService := api.Service{
		PlayTimeline:     backgroundService.PlayTimeline,
		StopTimeline:     backgroundService.StopTimeline,
		StopAllTimelines: backgroundService.StopAllTimelines,
		DB:               appdata,
		StartTime:        time.Now(),
	}

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

	//	Set up the API routes
	r := api.NewRouter(apiService)

	//	Start the timeline processor:
	go backgroundService.HandleAndProcess(ctx)

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
