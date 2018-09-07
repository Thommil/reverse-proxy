// Package users defines resources used for /users endpoint
package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/errors"

	"gitlab.forge.orange-labs.fr/mahali/services/ms-common/api"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-contributor/middlewares/authentication"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

type users struct {
	group       *gin.RouterGroup
	userService proto.UserService
}

// New creates new Routable implementation for /users resource
func New(engine *gin.Engine, userService proto.UserService, authenticationMiddleware gin.HandlerFunc) *gin.RouterGroup {
	users := &users{group: engine.Group("/users"), userService: userService}
	{
		//Public

		//Authenticated
		users.group.Use(authenticationMiddleware)
		users.group.GET("/:id", users.get)
		users.group.GET("", users.find)
		users.group.PUT("/:id", users.update)
	}
	return users.group
}

// @Summary Get a user from her ID
// @Description Get a user from her <b>id</b>. This API can be used to retrieve the current user's profile by using <i>me</i> parameter as <b>id</b>.
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "The user's ID, set to 'me' for current connected user"
// @Success 200 {object} proto.User
// @Failure 400 {object} errors.Error
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiJWT
// @Router /users/{id} [get]
func (users *users) get(c *gin.Context) {
	id := c.Param("id")
	if id == "me" {
		id = authentication.GetUser(c).Id
	}
	if user, err := users.userService.Get(context.TODO(), &proto.Id{Value: id}); err != nil {
		detailErr := errors.Parse(err.Error())
		c.AbortWithStatusJSON(int(detailErr.Code), detailErr)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// @Summary Get a list of users
// @Description <div>By default, without any parameters the full users paginated list is returned.<br/><br/>Search can be done using criterias directly in the query, in that case, a AND operator is applied:<pre>/users?role=admin&local=fr</pre><br/>For advanced search, the Mongo Query format is used, see <a href="https://docs.mongodb.com/manual/tutorial/query-documents/">here</a> for details:<pre>/users?query={"$or":[{"role":"admin},{"role":"user"}]}</pre></div>
// @Tags users
// @Accept  json
// @Produce  json
// @Param offset query int false "The result list start offset"
// @Param limit query int false "The result size limit"
// @Param sort query string false "The criteria used for sorting (ex: username), prefix with minus for reverse order (ex: -username)"
// @Param query query string false "The query for custom search (see description for details)"
// @Success 200 {object} api.Result
// @Failure 400 {object} errors.Error
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiJWT
// @Router /users [get]
func (users *users) find(c *gin.Context) {
	query, err := api.RequestToQuery(c.Request.URL)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"id":     "users_bad_request",
			"code":   http.StatusBadRequest,
			"detail": err.Error(),
			"status": http.StatusText(http.StatusBadRequest),
		})
	} else {
		result, err := users.userService.Find(context.TODO(), &proto.Query{
			Offset: int32(query.Offset),
			Limit:  int32(query.Limit),
			Sort:   query.Sort,
			Filter: query.Filter,
		})
		if err != nil {
			detailErr := errors.Parse(err.Error())
			c.AbortWithStatusJSON(int(detailErr.Code), detailErr)
		} else {
			if bodyResult, err := api.ResultToBody(result); err != nil {
				detailErr := errors.Parse(err.Error())
				c.AbortWithStatusJSON(int(detailErr.Code), detailErr)
			} else {
				c.JSON(http.StatusOK, bodyResult)
			}
		}
	}
}

// @Summary Update a user from her ID
// @Description Update a user from her <b>id</b>. Only current connected user can be updated if not admin.
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "The user's ID"
// @Success 200 {object} proto.User
// @Failure 400 {object} errors.Error
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiJWT
// @Router /users/{id} [put]
func (users *users) update(c *gin.Context) {
	id := c.Param("id")

	if id != authentication.GetUser(c).Id {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"id":     "users_forbidden",
			"code":   http.StatusForbidden,
			"detail": "cannot update another user",
			"status": http.StatusText(http.StatusForbidden),
		})
	} else {
		user := &proto.User{}
		if err := c.BindJSON(user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"id":     "users_bad_request",
				"code":   http.StatusBadRequest,
				"detail": err.Error(),
				"status": http.StatusText(http.StatusBadRequest),
			})
		} else {
			// Force admin only fields
			user.Id = id
			user.Role = ""

			if user, err = users.userService.Update(context.TODO(), user); err != nil {
				detailErr := errors.Parse(err.Error())
				c.AbortWithStatusJSON(int(detailErr.Code), detailErr)
			} else {
				c.JSON(http.StatusOK, user)
			}
		}
	}
}
