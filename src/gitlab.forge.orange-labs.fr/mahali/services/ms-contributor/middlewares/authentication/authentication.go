package authentication

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/micro/go-micro/errors"
	cache "github.com/patrickmn/go-cache"

	"github.com/gin-gonic/gin"
	authenticationProto "gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/proto"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

// Configuration defines configuration of Authentication
type Configuration struct {
	// Expired defines the delay in seconds before cache expiration
	Expired time.Duration `json:"expired"`
}

// Authenticated middleware check authentication
func Authenticated(configuration Configuration, authenticationService authenticationProto.AuthenticationService) gin.HandlerFunc {
	var userCache = cache.New(configuration.Expired*time.Second, configuration.Expired*time.Second)

	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("authorization")
		// First check presence of Authorization header
		if strings.Contains(authorizationHeader, "Bearer") {
			tokenString := strings.TrimSpace(strings.Replace(authorizationHeader, "Bearer", "", -1))
			if user, ok := userCache.Get(tokenString); ok {
				// Try to find user in local cache
				c.Set("user", user)
			} else {
				// // Try to get user from authentication service
				user, err := authenticationService.Validate(context.TODO(), &authenticationProto.Token{Value: tokenString})
				if err == nil {
					c.Set("user", user)
					userCache.Set(tokenString, user, cache.DefaultExpiration)
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, errors.Parse(err.Error()))
				}
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"id":     "authentication_unauthorized",
				"code":   http.StatusUnauthorized,
				"detail": "missing bearer",
				"status": http.StatusText(http.StatusUnauthorized),
			})
		}

		c.Next()
	}
}

// GetUser helper
func GetUser(c *gin.Context) *userProto.User {
	if user, found := c.Get("user"); found {
		return user.(*userProto.User)
	}
	return nil
}
