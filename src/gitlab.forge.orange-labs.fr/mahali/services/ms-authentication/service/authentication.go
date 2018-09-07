package service

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/micro/go-micro/errors"
	proto "gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/proto"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/providers"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

// JWTSettings defines JWT configuration
type JWTSettings struct {
	Secret  string        `json:"secret"`
	Expired time.Duration `json:"expired"`
	Issuer  string        `json:"issuer"`
}

type authenticationService struct {
	providers   map[string]providers.Provider
	jwtSettings JWTSettings
	userService userProto.UserService
}

// NewAuthenticationService creates a new instance of authentication service
func NewAuthenticationService(providersSettings map[string]map[string]string, session *mgo.Session, jwtSettings JWTSettings, userService userProto.UserService) (*authenticationService, error) {
	providersMap := make(map[string]providers.Provider)

	for name, settings := range providersSettings {
		provider, err := providers.NewInstance(name, settings, session, userService)
		if err != nil {
			return nil, err
		}
		providersMap[name] = provider
	}

	return &authenticationService{
		providers:   providersMap,
		jwtSettings: jwtSettings,
		userService: userService,
	}, nil
}

func (service *authenticationService) Authenticate(ctx context.Context, in *proto.AuthenticateRequest, out *proto.JWT) error {
	provider, found := service.providers[in.Provider]

	if !found {
		return errors.BadRequest("authentication_bad_request", "provider "+in.Provider+" not found")
	}

	user, err := provider.Authenticate(in.Credentials)
	if err != nil {
		detailErr := errors.Parse(err.Error())
		detailErr.Id = "authentication_unauthorized"
		return detailErr
	}

	claims := &jwt.StandardClaims{
		Subject:   user.Id,
		Issuer:    service.jwtSettings.Issuer,
		ExpiresAt: time.Now().Add(time.Second * service.jwtSettings.Expired).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(service.jwtSettings.Secret))
	if err != nil {
		return errors.Unauthorized("authentication_unauthorized", err.Error())
	}

	out.Token = ss
	out.ExpiresAt = claims.ExpiresAt

	return nil
}

func (service *authenticationService) Validate(ctx context.Context, in *proto.Token, out *userProto.User) error {
	// Check token headers
	token, err := jwt.Parse(in.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("authentication_unauthorized", "Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.jwtSettings.Secret), nil
	})

	if err != nil {
		return err
	}

	//Check token validity
	if token.Valid {
		//Check token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//Check user ID validity
			if userID, ok := claims["sub"]; ok {
				user, err := service.userService.Get(context.TODO(), &userProto.Id{Value: userID.(string)})
				if err != nil {
					detailErr := errors.Parse(err.Error())
					detailErr.Id = "authentication_unauthorized"
					return detailErr
				}

				//OK set user
				copier.Copy(out, user)

			} else {
				return errors.Unauthorized("authentication_unauthorized", "claim sub not found")
			}
		} else {
			return errors.Unauthorized("authentication_unauthorized", "invalid token")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errors.Unauthorized("authentication_unauthorized", "invalid token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return errors.Unauthorized("authentication_unauthorized", "token expired")
		} else {
			return errors.Unauthorized("authentication_unauthorized", err.Error())
		}
	} else {
		return errors.Unauthorized("authentication_unauthorized", err.Error())
	}

	return nil
}
