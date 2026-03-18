package context

// import (
// 	"github.com/zakiverse/zakiverse-api/core/cst"
// 	"github.com/gin-gonic/gin"
// )

// type AuthenticatedUser struct {
// 	ID   int64
// 	Type string
// }

// func GetAuthenticated(c *gin.Context) AuthenticatedUser {
// 	return AuthenticatedUser{
// 		ID:   c.GetInt64(cst.MiddlewareKeyTokenableId),
// 		Type: c.GetString(cst.MiddlewareKeyTokenableType),
// 	}
// }

// type Actor struct {
// 	ID   int64
// 	Type string
// }

// func GetActor(c *gin.Context) Actor {
// 	return Actor{
// 		ID:   c.GetInt64(cst.MiddlewareKeyActorId),
// 		Type: c.GetString(cst.MiddlewareKeyActorType),
// 	}
// }
