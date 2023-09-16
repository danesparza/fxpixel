package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAllTimelines godoc
// @Summary List all timelines in the system
// @Description List all timelines in the system
// @Tags timeline
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [get]
func (service Service) GetAllTimelines(rw http.ResponseWriter, req *http.Request) {

	//	Get a list of files
	dbTimelines, err := service.DB.GetAllTimelines(req.Context())
	if err != nil {
		err = fmt.Errorf("error getting a list of timelines: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	For each timeline, convert it to the API model:
	retval := []Timeline{}
	for _, timeline := range dbTimelines {
		apiTimeline := TimelineToApi(timeline)
		retval = append(retval, apiTimeline)
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("%v timelines(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}

// AddTimeline godoc
// @Summary Adds a timeline to the system
// @Description Adds a timeline to the system
// @Tags timeline
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [put]
func (service Service) AddTimeline(rw http.ResponseWriter, req *http.Request) {

	//	Parse the body

	//	Get a list of files
	dbTimelines, err := service.DB.GetAllTimelines(req.Context())
	if err != nil {
		err = fmt.Errorf("error getting a list of timelines: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	For each timeline, convert it to the API model:
	retval := []Timeline{}
	for _, timeline := range dbTimelines {
		apiTimeline := TimelineToApi(timeline)
		retval = append(retval, apiTimeline)
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("%v timelines(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}
