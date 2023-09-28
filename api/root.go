package api

import (
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/leds"

	"github.com/danesparza/fxpixel/version"
	"net/http"
	"time"
)

// Service encapsulates API service operations
type Service struct {
	// DB is the system datastore reference
	DB data.AppDataService

	// StartTime is the service start time.  We can calculate uptime based on this
	StartTime time.Time

	// PlayTimeline signals a timeline should be played
	PlayTimeline chan leds.PlayTimelineRequest

	// StopTimeline signals a timeline should stop playing
	StopTimeline chan string

	//	StopAllTimelines signals all timelines should stop playing
	StopAllTimelines chan bool
}

// UpdateTagsRequest represents a request to update tags for a file
type UpdateTagsRequest struct {
	Tags []string `json:"tags"`
}

// UpdateConfigValueRequest represents a request to update a system config value
type UpdateConfigValueRequest struct {
	Value int `json:"value"`
}

// SystemResponse is a response for a system request
type SystemResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an API response
type ErrorResponse struct {
	Message string `json:"message"`
}

// Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := ErrorResponse{
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(response)
}

// ShowUI redirects to the /ui/ url path
func (service Service) ShowUI(rw http.ResponseWriter, req *http.Request) {
	// http.Redirect(rw, req, "/ui/", 301)
	fmt.Fprintf(rw, "Hello, world - UI")
}

// ApiVersionMiddleware adds the API version information to the response header
func ApiVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//	Include the version in the response headers:
		rw.Header().Set(version.Header, version.String())

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
