package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	adding "Project/adding"
	"Project/database/dynamodb"
	redisdb "Project/database/redis"
	apiGatewayService "Project/handlers/apigatewayService"
	rest "Project/handlers/rest"
	listing "Project/listing"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbType := flag.String("database", "dynamo", "database type [redis, dynamo]")
	handlerType := flag.String("handler", "create", "handler type [localhost, listall create]")
	flag.Parse()

	var addingRepo adding.DroneRepository
	var listingRepo listing.DroneRepository

	switch *dbType {
	case "redis":
		addingRepo = redisdb.NewRedisAddingRepository(redisConnection("localhost:6379"))
		listingRepo = redisdb.NewRedisListingRepository(redisConnection("localhost:6379"))
	case "dynamo":
		addingRepo = dynamodb.NewDynamoAddingRepository(dynamoDBConnection())
		listingRepo = dynamodb.NewDynamoListingRepository(dynamoDBConnection())

	default:
		panic("Unknown database")
	}

	droneAdding := adding.NewDroneService(addingRepo)
	droneListing := listing.NewDroneService(listingRepo)

	switch *handlerType {
	case "localhost":
		fmt.Println("localhost handler")
		addingHandler := rest.NewAddingHandler(droneAdding)
		listingHandler := rest.NewListingHandler(droneListing)
		serverLocalhost(addingHandler, listingHandler)
	case "listall":
		fmt.Println("AWS API Gateway handler Listall")
		droneApiGateway := apiGatewayService.NewListingApiGateway(droneListing)
		apiGatewayGetList(droneApiGateway)
	case "create":
		fmt.Println("AWS API Gateway handler CreateDrone")
		droneApiGateway := apiGatewayService.NewAddingApiGateway(droneAdding)
		apiGatewayCreateDrone(droneApiGateway)
	default:
		panic("Unknown handler")
	}
}

func apiGatewayGetList(droneApiGateway apiGatewayService.ListingApiGateway) {
	lambda.Start(droneApiGateway.Get)
}

func apiGatewayCreateDrone(droneApiGateway apiGatewayService.AddingApiGateway) {
	lambda.Start(droneApiGateway.Create)
}

func serverLocalhost(addingHandler rest.AddingHandler, listingHandler rest.ListingHandler) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/drones", listingHandler.Get).Methods("GET")
	router.HandleFunc("/drones", addingHandler.Create).Methods("POST")

	errs := make(chan error, 2)

	go func() {
		http.Handle("/", accessControl(router))
		fmt.Println("Listening on port :3001")
		errs <- http.ListenAndServe(":3001", nil)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("terminated %s", <-errs)
}

func redisConnection(url string) *redis.Client {
	fmt.Println("Connecting to Redis DB")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})
	err := client.Ping().Err()

	if err != nil {
		panic(err)
	}
	return client
}

func dynamoDBConnection() *session.Session {
	fmt.Println("Connecting to DynamoDB DB")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
