package rest

import (
	"Project/listing"
	"encoding/json"
	"net/http"
)

type ListingHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type droneListingHandler struct {
	droneService listing.DroneService
}

func NewListingHandler(droneService listing.DroneService) ListingHandler {
	return &droneListingHandler{
		droneService,
	}
}

func (h *droneListingHandler) Get(w http.ResponseWriter, r *http.Request) {
	drones, _ := h.droneService.FindAllDrones()

	response, _ := json.Marshal(drones)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
