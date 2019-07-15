package adding

import (
	"github.com/google/uuid"
)

type DroneRepository interface {
	Create(drone *Drone) error
}

type DroneService interface {
	CreateDrone(drone *Drone) error
}

type droneService struct {
	repo DroneRepository
}

func NewDroneService(repo DroneRepository) DroneService {
	return &droneService{
		repo,
	}
}

func (s *droneService) CreateDrone(drone *Drone) error {
	drone.ID = uuid.New().String()
	return s.repo.Create(drone)
}
