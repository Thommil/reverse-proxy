package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/globalsign/mgo"
	micro "github.com/micro/go-micro"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-user/service"
)

// Configuration definition for ms-user
type Configuration struct {
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
		micro.Name("user"),
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
	if err = db.DB("").C("configuration").FindId("user").One(config); err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	// Register services
	err = proto.RegisterUserServiceHandler(ms.Server(), service.NewUserService(db))
	if err != nil {
		log.Fatalf("Failed to register user service : %v\n", err)
	}

	// Start Server
	if err = ms.Run(); err != nil {
		log.Fatalf("Failed to start micro : %v\n", err)
	}
}
