// Package authentication defines resources used for /authentication endpoint
package authentication

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/errors"

	"gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/proto"
)

type authentication struct {
	group                 *gin.RouterGroup
	authenticationService proto.AuthenticationService
}

// New creates new Routable implementation for /users resource
func New(engine *gin.Engine, authenticationService proto.AuthenticationService, authenticationMiddleware gin.HandlerFunc) *gin.RouterGroup {
	authentication := &authentication{group: engine.Group("/authentication"), authenticationService: authenticationService}
	{
		//Public
		authentication.group.POST("/:provider", authentication.authenticate)
	}
	return authentication.group
}

// @Summary Authenticate client based on given provider and credentials
// @Description <div>Creates a new JWT token to be used in headers as follow:<br/><pre>Authorization : Bearer {JWT_Token}</pre><br/>The authentication <b>provider</b> is set directly in path.   The model of <b>credentials</b> sent in body is based on selected <b>provider</b>:<ul><li><i>local</i> : uses the localy stored credentials in Mahali service:</li></ul><pre>{<br/>&nbsp;&nbsp;"username" : {username},<br/>&nbsp;&nbsp;"password" : {password}<br/>}</pre></div>
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param provider path string true "Provider used for authentication (ex : local, google, facebook ...)"
// @Param credentials body interface{} true "Depends on provider, see description above for details"
// @Success 200 {object} proto.JWT
// @Failure 400 {object} errors.Error
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /authentication/{provider} [post]
func (authentication *authentication) authenticate(c *gin.Context) {
	provider := c.Param("provider")
	credentials := make(map[string]string)

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"id":     "authentication_bad_request",
			"code":   http.StatusBadRequest,
			"detail": "invalid credentials",
			"status": http.StatusText(http.StatusBadRequest),
		})
	} else {
		if jwt, err := authentication.authenticationService.Authenticate(context.TODO(), &proto.AuthenticateRequest{
			Provider:    provider,
			Credentials: credentials,
		}); err != nil {
			jsonErr := errors.Parse(err.Error())
			c.AbortWithStatusJSON(int(jsonErr.Code), jsonErr)
		} else {
			c.JSON(http.StatusOK, jwt)
		}
	}
}
