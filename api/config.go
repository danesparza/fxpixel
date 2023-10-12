package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSystemConfig godoc
// @Summary Get the system configuration information
// @Description Get the system configuration information
// @Tags timeline
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /config [get]
func (service Service) GetSystemConfig(rw http.ResponseWriter, req *http.Request) {

	//	Get a list of files
	systemConfig, err := service.DB.GetSystemConfig(req.Context())
	if err != nil {
		err = fmt.Errorf("error getting system config: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: "Fetched system config",
		Data:    systemConfig,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}
