package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/globalsign/mgo"
	micro "github.com/micro/go-micro"
	proto "gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/proto"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/service"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

// Configuration definition for ms-authentication
type Configuration struct {
	// JWT settings from service
	JWT service.JWTSettings

	// Providers for authentication
	Providers map[string]map[string]string
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

func main() {
	// Instanciate Micro service
	ms := micro.NewService(
		micro.Name("authentication"),
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
	if err = db.DB("").C("configuration").FindId("authentication").One(config); err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	// Register services
	handler, err := service.NewAuthenticationService(config.Providers, db, config.JWT, userProto.NewUserService("user", ms.Client()))
	if err != nil {
		log.Fatalf("Failed to instanciate authentication service : %v\n", err)
	}

	err = proto.RegisterAuthenticationServiceHandler(ms.Server(), handler)
	if err != nil {
		log.Fatalf("Failed to register authentication service : %v\n", err)
	}

	// Start Server
	if err = ms.Run(); err != nil {
		log.Fatalf("Failed to start micro : %v\n", err)
	}
}
