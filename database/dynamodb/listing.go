package dynamodb

import (
	"Project/listing"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type droneListingRepository struct {
	sess *session.Session
}

func NewDynamoListingRepository(sess *session.Session) listing.DroneRepository {
	return &droneListingRepository{
		sess,
	}
}

func (r *droneListingRepository) FindAll() (drones []*listing.Drone, err error) {
	svc := dynamodb.New(r.sess)
	input := &dynamodb.ScanInput{
		TableName: aws.String("drones"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	var droneList []*listing.Drone
	for _, element := range result.Items {
		var dron listing.Drone
		dron.ID = element["ID"].String()
		dron.Name = element["Name"].String()
		dron.Description = element["Description"].String()
		droneList = append(droneList, &dron)
	}
	return droneList, nil

}
