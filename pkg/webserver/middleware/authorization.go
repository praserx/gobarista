package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// reqPerm := authorization.GetPermissionFromRequestMethod(c.Request.Method)
		// if !authorization.IsAuthorized("/", "admin", reqPerm, false) {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		// 	return
		// }

		c.Next()
	}
}
