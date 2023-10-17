package leds

import (
	"context"
	"github.com/Jon-Bright/ledctl/pixarray"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	stepType "github.com/danesparza/fxpixel/internal/data/const/step"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type PlayTimelineRequest struct {
	ProcessID         string
	RequestedTimeline data.Timeline
}

type timelineProcessMap struct {
	m       map[string]func()
	rwMutex sync.RWMutex
}

// StepProcessor encapsulates the core config for processing a step
type StepProcessor struct {
	GPIO           int
	LEDs           int
	PixelOrder     string
	NumberOfColors int
	PixArray       *pixarray.PixArray
}

// BackgroundProcess encapsulates background processing operations
type BackgroundProcess struct {
	// DB is the system datastore reference
	DB data.AppDataService

	// PlayTimeline signals a timeline should be played
	PlayTimeline chan PlayTimelineRequest

	// StopTimeline signals a running timeline should be stopped
	StopTimeline chan string

	// StopAllTimelines signals all running timlines should be stopped
	StopAllTimelines chan bool

	// PlayingTimelines tracks currently playing timelines
	PlayingTimelines timelineProcessMap
}

// HandleAndProcess handles system context calls and channel events to play/stop timelines
func (bp *BackgroundProcess) HandleAndProcess(systemctx context.Context) {

	//	Create a map of running timelines and their cancel functions
	bp.PlayingTimelines.m = make(map[string]func())
	log.Info().Msg("Starting timeline processor...")

	//	Loop and respond to channels:
	for {
		select {
		case playReq := <-bp.PlayTimeline:
			//	As we get a request on a channel to play a file...
			//	Spawn a goroutine
			log.Debug().Str("TimelineId", playReq.RequestedTimeline.ID).Msg("Starting to play timeline")
			go bp.StartTimelinePlay(systemctx, playReq) // Launch the goroutine

		case stopTL := <-bp.StopTimeline:

			//	Look up the item in the map and call cancel if the item exists (critical section):
			bp.PlayingTimelines.rwMutex.Lock()
			playCancel, exists := bp.PlayingTimelines.m[stopTL]

			if exists {
				//	Call the context cancellation function
				playCancel()

				log.Debug().Str("ProcessID", stopTL).Msg("Stopped timeline process")

				//	Remove ourselves from the map and exit
				delete(bp.PlayingTimelines.m, stopTL)
			}
			bp.PlayingTimelines.rwMutex.Unlock()

		case <-bp.StopAllTimelines:

			//	Loop through all items in the map and call cancel if the item exists (critical section):
			bp.PlayingTimelines.rwMutex.Lock()

			log.Debug().Msg("Stopping all timeline processes")

			for stopTL, playCancel := range bp.PlayingTimelines.m {

				//	Call the cancel function
				playCancel()

				//	Remove ourselves from the map
				//	(this is safe to do in a 'range':
				//	https://golang.org/doc/effective_go#for )
				delete(bp.PlayingTimelines.m, stopTL)
			}

			bp.PlayingTimelines.rwMutex.Unlock()

		case <-systemctx.Done():
			log.Info().Msg("Stopping timeline processor")
			return
		}
	}
}

// PlayTimeline plays a timeline
func (bp *BackgroundProcess) StartTimelinePlay(cx context.Context, req PlayTimelineRequest) {
	//	Create a cancelable context from the passed context
	ctx, cancel := context.WithCancel(cx)
	defer cancel()

	//	Add an entry to the map with
	//	- key: instance id
	//	- value: the cancel function (pointer)
	//	(critical section)
	bp.PlayingTimelines.rwMutex.Lock()
	bp.PlayingTimelines.m[req.ProcessID] = cancel
	bp.PlayingTimelines.rwMutex.Unlock()

	//	Get the system default configuration
	systemConfig, err := bp.DB.GetSystemConfig(ctx)
	if err != nil {
		log.Err(err).Msg("An error occurred trying to get the system config")
		return
	}

	//	Spin up a strip:
	pixels, err := NewStrip( // Take the defaults for most things ...
		systemConfig.LEDs,                               // Set the number of LEDs
		WithGPIOPIn(systemConfig.GPIO),                  // Set the GPIO pin
		WithPixelOrder(systemConfig.PixelOrder),         // Set the pixel order
		WithNumberOfColors(systemConfig.NumberOfColors), // Set the number of colors
	)
	if err != nil {
		log.Err(err).Msg("Problem creating strip")
		return
	}

	//	Create a new pixel array
	arr := pixarray.NewPixArray(systemConfig.LEDs, systemConfig.NumberOfColors, pixels)

	//	Set the defaults for the StepProcessor:
	sp := StepProcessor{
		GPIO:           systemConfig.GPIO,
		LEDs:           systemConfig.LEDs,
		PixelOrder:     systemConfig.PixelOrder,
		NumberOfColors: systemConfig.NumberOfColors,
		PixArray:       arr,
	}

	//	Process the timeline
	log.Debug().Str("ProcessID", req.ProcessID).Msg("Processing timeline")

	//	First, see if the timeline has a GPIO port set on it.
	if req.RequestedTimeline.GPIO.Valid == true || req.RequestedTimeline.GPIO.Int32 != 0 {
		//	If so, use that information:
		sp.GPIO = int(req.RequestedTimeline.GPIO.Int32)
	}

	//	Keep a channel state map:
	//channelState := map[int]byte{}

	//	Our waitgroup (for sync'ing fade finishes)
	//var wg sync.WaitGroup

	//	Iterate through each step
	for _, step := range req.RequestedTimeline.Steps {

		select {
		default:

			//	Find out what type of frame this is, and act accordingly:
			switch step.Type {
			case stepType.Unknown:
				//	We're not sure what happened, but this can't be processed.
				log.Warn().
					Str("timelineid", req.RequestedTimeline.ID).
					Str("stepid", step.ID).
					Msg("Step has unknown steptype and can't be processed")

			case stepType.Loop:
				//	Get the loop information and process the loop:
				log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing loop")

			case stepType.Trigger:
				//	Get the trigger information and process the trigger:
				log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing trigger")

			case stepType.Sleep:
				//	Get the sleep information and pause here
				log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing sleep")

			case stepType.RandomSleep:
				//	Get the random sleep parameters and pause here
				log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing random sleep")

			case stepType.Effect:
				//	Find the effect type and process it.
				switch step.Effect {
				case effect.Fade:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: fade")

				case effect.Gradient:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: gradient")

				case effect.KnightRider:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: knight rider")

				case effect.Lightning:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: lightning")

				case effect.Rainbow:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: rainbow")

				case effect.Sequence:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: sequence")

				case effect.Solid:
					sp.ProcessSolidEffect(step)

					//	Sleep for the time specified
					//	(this has the effect of showing the color for this amount of time)
					select {
					case <-time.After(time.Duration(step.Time.Int32) * time.Millisecond):
						continue
					case <-ctx.Done():
						return
					}

				case effect.Zip:
					log.Debug().Str("stepid", step.ID).Int32("time", step.Time.Int32).Msg("Processing effect: zip")

				}

			}

		case <-ctx.Done():
			// stop
			return
		}
	}

	//	Remove ourselves from the map and exit (critical section)
	bp.PlayingTimelines.rwMutex.Lock()
	delete(bp.PlayingTimelines.m, req.ProcessID)
	bp.PlayingTimelines.rwMutex.Unlock()

	log.Debug().Str("ProcessID", req.ProcessID).Msg("Processing completed for timeline")
}
