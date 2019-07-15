package rest

import (
	"Project/adding"
	"encoding/json"
	"net/http"
)

type AddingHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type droneAddingHandler struct {
	droneService adding.DroneService
}

func NewAddingHandler(droneService adding.DroneService) AddingHandler {
	return &droneAddingHandler{
		droneService,
	}
}

func (h *droneAddingHandler) Create(w http.ResponseWriter, r *http.Request) {

	var drone adding.Drone
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&drone)
	_ = h.droneService.CreateDrone(&drone)

	response, _ := json.Marshal(drone)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}
