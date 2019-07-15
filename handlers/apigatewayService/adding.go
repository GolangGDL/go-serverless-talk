package apigatewayService

import (
	"Project/adding"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

type AddingApiGateway interface {
	Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type droneAddingApiGateway struct {
	droneService adding.DroneService
}

func NewAddingApiGateway(droneService adding.DroneService) AddingApiGateway {
	return &droneAddingApiGateway{
		droneService,
	}
}

func (h *droneAddingApiGateway) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var drone adding.Drone
	
	err := json.Unmarshal([]byte(request.Body), &drone)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid payload",
		}, nil
	}
	
	_ = h.droneService.CreateDrone(&drone)

	response, _ := json.Marshal(drone)

	_, err = external.LoadDefaultAWSConfig()
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
