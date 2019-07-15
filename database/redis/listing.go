package redis

import (
	"Project/listing"
	"encoding/json"

	"github.com/go-redis/redis"
)

//const table = "drones"

type droneListingRepository struct {
	connection *redis.Client
}

func NewRedisListingRepository(connection *redis.Client) listing.DroneRepository {
	return &droneListingRepository{
		connection,
	}
}

func (r *droneListingRepository) FindAll() (drones []*listing.Drone, err error) {
	const table = "drones"

	ts := r.connection.HGetAll(table).Val()
	for key, value := range ts {
		t := new(listing.Drone)
		err = json.Unmarshal([]byte(value), t)

		if err != nil {
			return nil, err
		}

		t.ID = key
		drones = append(drones, t)
	}
	return drones, nil
}
