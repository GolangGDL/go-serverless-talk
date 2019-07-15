package apigatewayService

import (
	"Project/listing"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

type ListingApiGateway interface {
	Get() (events.APIGatewayProxyResponse, error)
}

type droneListingApiGateway struct {
	droneService listing.DroneService
}

func NewListingApiGateway(droneService listing.DroneService) ListingApiGateway {
	return &droneListingApiGateway{
		droneService,
	}
}

func (h *droneListingApiGateway) Get() (events.APIGatewayProxyResponse, error) {

	drones, _ := h.droneService.FindAllDrones()
	response, _ := json.Marshal(drones)

	_, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, err
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}
