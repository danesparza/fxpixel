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

// GetAllTimelinesWithTag godoc
// @Summary Gets timelines that have a tag
// @Description Gets timelines that have a tag
// @Tags timeline
// @Accept  json
// @Produce  json
// @Param tag path string true "The tag to use when fetching timelines"
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines/tag/{tag} [get]
func (service Service) GetAllTimelinesWithTag(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	tag := chi.URLParam(req, "tag")

	//	Add a timeline
	dbTimelines, err := service.DB.GetAllTimelinesWithTag(req.Context(), tag)
	if err != nil {
		err = fmt.Errorf("error getting a timeline: %v", err)
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
// @Param endpoint body api.Timeline true "The timeline to add"
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
// @Param id path string true "The timeline id to get"
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines/{id} [get]
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
// @Param id path string true "The timeline id to delete"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines/{id} [delete]
func (service Service) DeleteTimeline(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	timelineId := chi.URLParam(req, "id")

	if timelineId == "" {
		err := fmt.Errorf("requires an id of a timeline to delete")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

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

// UpdateTags godoc
// @Summary Updates tags for a timeline
// @Description Updates tags for a timeline
// @Tags timeline
// @Accept  json
// @Produce  json
// @Param id path string true "The timeline id to update tags for"
// @Param endpoint body api.UpdateTagsRequest true "The tags to set for the timeline"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /timelines/{id} [post]
func (service Service) UpdateTags(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	timelineId := chi.URLParam(req, "id")

	if timelineId == "" {
		err := fmt.Errorf("requires an id of a file to update tags for")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Parse the body to get the tags
	request := UpdateTagsRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		err = fmt.Errorf("problem decoding tag update request: %v", err)
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Update the timeline
	err = service.DB.UpdateTags(req.Context(), timelineId, request.Tags)
	if err != nil {
		err = fmt.Errorf("error updating tags: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: "Timeline updated",
		Data:    timelineId,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}
