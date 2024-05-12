package pages

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/webserver/security"
)

const RootTitle = "Main | GoBarista"

func RootGET(c *gin.Context) {
	c.HTML(http.StatusOK, "root.go.tmpl", gin.H{
		"title": "home",
	})
}

func RootPOST(c *gin.Context) {
	var ok bool
	var email string

	if email, ok = c.GetPostForm("email"); !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	logger.Info(email)

	if user, err := database.SelectUserByEmail(email); err == nil {

		sessionID := sessions.Default(c).Get(security.SessionKey).(string)

		userSession := security.Session{
			UserID:      user.ID,
			UserRole:    "",
			Code:        security.Code(),
			CodeValidTo: time.Now().Add(5 * time.Minute),
			Logged:      false,
		}

		security.SessionSet(sessionID, userSession)
		logger.Info(fmt.Sprintf("code: %s", userSession.Code))
		// Send verification code here
	} else {
		logger.Info("user do not exists")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Redirect(http.StatusSeeOther, "/code-verification")
}
