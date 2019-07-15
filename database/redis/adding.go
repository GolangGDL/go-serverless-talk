package redis

import (
	"Project/adding"
	"encoding/json"

	"github.com/go-redis/redis"
)

const table = "drones"

type droneAddingRepository struct {
	connection *redis.Client
}

func NewRedisAddingRepository(connection *redis.Client) adding.DroneRepository {
	return &droneAddingRepository{
		connection,
	}
}

func (r *droneAddingRepository) Create(drone *adding.Drone) error {
	encoded, err := json.Marshal(drone)

	if err != nil {
		return err
	}

	r.connection.HSet(table, drone.ID, encoded)
	return nil
}
