package google

import (
	"github.com/globalsign/mgo"

	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

// Allows to check user entry against Google account
type provider struct {
	configuration map[string]string
	db            *mgo.Session
	userService   *userProto.UserService
}

// New instanciates a new Google provider
func New(configuration map[string]string, db *mgo.Session, userService *userProto.UserService) (*provider, error) {
	return &provider{
		configuration: configuration,
		db:            db,
		userService:   userService,
	}, nil
}

func (p provider) Authenticate(credentials map[string]string) (*userProto.User, error) {
	return nil, nil
}

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/thommil/ms-common/model"
// )

// // Configuration definition for google provider
// type Configuration struct {
// 	// URL used to check token (:idToken replaced)
// 	URL string `json:"url"`
// 	// Issuer to check
// 	ISS string `json:"iss"`
// 	// Audience to check
// 	AUD string `json:"aud"`
// }

// type tokenInfo struct {
// 	AZP        string `json:"azp"`
// 	AUD        string `json:"aud"`
// 	SUB        string `json:"sub"`
// 	EXP        string `json:"exp"`
// 	ISS        string `json:"iss"`
// 	IAT        string `json:"iat"`
// 	Name       string `json:"name"`
// 	Picture    string `json:"picture"`
// 	GivenName  string `json:"given_name"`
// 	FamilyName string `json:"family_name"`
// 	Locale     string `json:"locale"`
// 	ALG        string `json:"alg"`
// 	KID        string `json:"kid"`
// }

// // Provider allows to check user entry against Google JWT token
// type Provider struct {
// 	Configuration *Configuration
// }

// // Authenticate implementation of Provider API
// func (provider Provider) Authenticate(credentials interface{}) (*model.User, *model.Error) {
// 	var httpClient = &http.Client{Timeout: 10 * time.Second}

// 	response, err := httpClient.Get(strings.Replace(provider.Configuration.URL, ":idToken", credentials.(string), 1))
// 	if err != nil {
// 		return nil, &model.Error{
// 			Code:    http.StatusInternalServerError,
// 			Message: err.Error(),
// 		}
// 	} else if response.StatusCode >= http.StatusBadRequest {
// 		return nil, &model.Error{
// 			Code:    response.StatusCode,
// 			Message: response.Status,
// 		}
// 	}

// 	defer response.Body.Close()
// 	token := &tokenInfo{}
// 	json.NewDecoder(response.Body).Decode(token)
// 	if strings.Contains(token.ISS, provider.Configuration.ISS) && token.AUD == provider.Configuration.AUD {
// 		query := model.FindAccount(&model.Account{ExternalID: token.SUB})
// 		count, err := query.Count()
// 		if err != nil {
// 			return nil, &model.Error{
// 				Code:    http.StatusInternalServerError,
// 				Message: err.Error(),
// 			}
// 		}
// 		if count > 0 {
// 			//Found return user
// 			account := &model.Account{}
// 			if query.One(account) != nil {
// 				return nil, &model.Error{
// 					Code:    http.StatusInternalServerError,
// 					Message: err.Error(),
// 				}
// 			}
// 			//return model.FindUserByID(account.UserID)
// 		}

// 		//Not found, create account & user
// 		/*user, err := model.CreateOrUpdateUser(&model.User{Username: token.Name, Picture: token.Picture, Locale: token.Locale})
// 		if err != nil {
// 			return nil, err
// 		}
// 		model.CreateOrUpdateAccount(&model.Account{Provider: "google", ExternalID: token.SUB, UserID: user.ID.Hex()})
// 		return user, nil*/
// 	}

// 	return nil, &model.Error{
// 		Code:    http.StatusBadRequest,
// 		Message: "Bad 'iss' or 'aud' claim in Google token",
// 	}
// }
