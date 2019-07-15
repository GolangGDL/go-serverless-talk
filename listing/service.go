package listing

type DroneRepository interface {
	FindAll() ([]*Drone, error)
}

type DroneService interface {
	FindAllDrones() ([]*Drone, error)
}

type droneService struct {
	repo DroneRepository
}

func NewDroneService(repo DroneRepository) DroneService {
	return &droneService{
		repo,
	}
}

func (s *droneService) FindAllDrones() ([]*Drone, error) {
	return s.repo.FindAll()
}
