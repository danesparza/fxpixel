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
