package providers

import (
	"fmt"

	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"

	"github.com/globalsign/mgo"

	local "gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/providers/local"
)

// Provider defines authentication providers interface
type Provider interface {
	// Authenticate a user given the credentials
	Authenticate(credentials map[string]string) (*userProto.User, error)
}

// NewInstance instanciates a new instance of a Provider from its name
func NewInstance(name string, configuration map[string]string, db *mgo.Session, userService userProto.UserService) (Provider, error) {
	// Refect is tempting but totaly overkill here
	if name == "local" {
		return local.New(configuration, db, userService)
	} else {
		return nil, fmt.Errorf("provider not found %s", name)
	}
}
