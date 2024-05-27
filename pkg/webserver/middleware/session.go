package middleware

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/webserver/security"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get(security.SessionKey) == nil {
			sk, err := uuid.NewRandom()
			if err != nil {
				logger.Error(fmt.Sprintf("cannot initialize session properly: cannot create random uuid: %v", err))
			} else {
				session.Set(security.SessionKey, sk.String())
				session.Save()
			}
		}

		c.Next()
	}
}
