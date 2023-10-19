package api

import (
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxpixel/internal/leds"
	"github.com/go-chi/chi/v5"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"net/http"
)

// RequestTimelinePlay godoc
// @Summary Plays a timeline in the system
// @Description Plays a timeline in the system
// @Tags process
// @Accept  json
// @Produce  json
// @Param id path string true "The timeline id to play"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timeline/run/{id} [post]
func (service Service) RequestTimelinePlay(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	timelineid := chi.URLParam(req, "id")

	//	Get the timeline
	timeline, err := service.DB.GetTimeline(req.Context(), timelineid)
	if err != nil {
		err = fmt.Errorf("error getting timeline: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	log.Debug().Str("timelineid", timeline.ID).Str("timeline_name", timeline.Name).Msg("Requesting timeline play")

	//	Send to the channel:
	playRequest := leds.PlayTimelineRequest{
		ProcessID:         xid.New().String(), // Generate a new id
		RequestedTimeline: timeline,
	}
	service.PlayTimeline <- playRequest

	//	Construct our response
	response := SystemResponse{
		Message: "Requested timeline play",
		Data:    playRequest,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// RequestTimelineStop godoc
// @Summary Stops a specific timeline 'play' process
// @Description Stops a specific timeline 'play' process
// @Tags process
// @Accept  json
// @Produce  json
// @Param pid path string true "The process id to stop"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /timeline/stop/{pid} [post]
func (service Service) RequestTimelineStop(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	timelinepid := chi.URLParam(req, "pid")
	if timelinepid == "" {
		err := fmt.Errorf("requires a processid of a process to stop")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Send to the channel:
	service.StopTimeline <- timelinepid

	log.Debug().Str("pid", timelinepid).Msg("Requesting timeline process stop")

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Timeline stopping",
		Data:    timelinepid,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// RequestAllTimelinesStop godoc
// @Summary Stops all timeline 'play' processes
// @Description Stops all timeline 'play' processes
// @Tags process
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Router /timeline/stop [post]
func (service Service) RequestAllTimelinesStop(rw http.ResponseWriter, req *http.Request) {

	//	Send to the channel:
	service.StopAllTimelines <- true

	log.Debug().Msg("Requesting all timeline processes stop")

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "All Timelines stopping",
		Data:    ".",
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
