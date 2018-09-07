package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	micro "github.com/micro/go-micro"
	authenticationProto "gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/proto"
	_ "gitlab.forge.orange-labs.fr/mahali/services/ms-contributor/docs"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-contributor/middlewares/authentication"
	authenticationResource "gitlab.forge.orange-labs.fr/mahali/services/ms-contributor/resources/authentication"
	usersResource "gitlab.forge.orange-labs.fr/mahali/services/ms-contributor/resources/users"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"

	"github.com/swaggo/gin-swagger"              // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
)

// Configuration definition for ms-authentication
type Configuration struct {
	HTTP struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	Authentication authentication.Configuration `json:"authentication"`
}

//TODO Improve DB discovery (cluster support & DB name)
func mongoSession(ms micro.Service) (*mgo.Session, error) {
	services, err := ms.Options().Registry.GetService("mongo")

	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no mongo service found in registry")
	}

	session, err := mgo.Dial("mongodb://" + services[0].Nodes[0].Address + ":" + strconv.Itoa(services[0].Nodes[0].Port) + "/test")

	if err != nil {
		return nil, err
	}

	return session, nil
}

// @title Mahali API
// @version 1.0
// @description Mahali Service API definition

// @securityDefinitions.apikey ApiJWT
// @in header
// @name Authorization
func main() {
	// Instanciate Micro service
	ms := micro.NewService(
		micro.Name("contributor"),
	)
	ms.Init()

	// Init Mongo Session
	db, err := mongoSession(ms)
	if err != nil {
		log.Fatalf("Failed to connect to database : %v\n", err)
	}
	defer db.Close()

	// Configuration
	config := &Configuration{}
	if err = db.DB("").C("configuration").FindId("contributor").One(config); err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	// Services clients
	userService := userProto.NewUserService("user", ms.Client())
	authenticationService := authenticationProto.NewAuthenticationService("authentication", ms.Client())

	//GIN Server
	router := gin.Default()

	//Middlewares
	authenticationMiddleware := authentication.Authenticated(config.Authentication, authenticationService) // Authent
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))                              // Swagger

	//Resources
	authenticationResource.New(router, authenticationService, authenticationMiddleware)
	usersResource.New(router, userService, authenticationMiddleware)

	//Start GIN Server in Go Routine
	go endless.ListenAndServe(fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port), router)

	// Start Server
	if err = ms.Run(); err != nil {
		log.Fatalf("Failed to start micro : %v\n", err)
	}
}
