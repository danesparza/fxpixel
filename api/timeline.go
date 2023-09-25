package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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
	request := Timeline{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		err = fmt.Errorf("problem decoding add timeline request: %v", err)
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Convert the api request into a data model:
	newTimeline := ApiToTimeline(request)

	//	Add a timeline
	dbTimeline, err := service.DB.AddTimeline(req.Context(), newTimeline)
	if err != nil {
		err = fmt.Errorf("error adding a timelines: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Convert the timeline to the API model:
	retval := TimelineToApi(dbTimeline)

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("Timeline added: %v", retval.ID),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}

// GetTimeline godoc
// @Summary Gets a single timeline
// @Description Gets a single timeline
// @Tags timeline
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [put]
func (service Service) GetTimeline(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	timelineId := chi.URLParam(req, "id")

	//	Add a timeline
	dbTimeline, err := service.DB.GetTimeline(req.Context(), timelineId)
	if err != nil {
		err = fmt.Errorf("error getting a timeline: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Convert the timeline to the API model:
	retval := TimelineToApi(dbTimeline)

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("Timeline fetched: %v", retval.ID),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}

// DeleteTimeline godoc
// @Summary Delete a single timeline
// @Description Delete a single timeline
// @Tags timeline
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines [delete]
func (service Service) DeleteTimeline(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	timelineId := chi.URLParam(req, "id")

	//	Add a timeline
	err := service.DB.DeleteTimeline(req.Context(), timelineId)
	if err != nil {
		err = fmt.Errorf("error deleting a timeline: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("Timeline deleted: %v", timelineId),
		Data:    timelineId,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}
